package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func makeRequest(url string, method string) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request object
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	apiKey := os.Getenv("OPENAI_API_SECRET_KEY")
	req.Header.Set("accept", "*/*")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		bodyText, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyText))
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	return resp, nil
}

func GetOrgUsers(orgId string) (*UsersResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/organizations/%s/users", orgId)
	resp, err := makeRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var usersResponse UsersResponse
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return nil, err
	}
	return &usersResponse, err
}

func GetDayUsage(user User, date string) (*DailyUsageResponse, error) {
	usageRequestUrl := fmt.Sprintf(
		"https://api.openai.com/v1/usage?date=%s&user_public_id=%s",
		date,
		user.ID,
	)
	resp, err := makeRequest(usageRequestUrl, http.MethodGet)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var dailyUsageResp DailyUsageResponse
	err = json.Unmarshal(body, &dailyUsageResp)
	if err != nil {
		return nil, err
	}
	return &dailyUsageResp, err
}

func CalculateUserUsage(user User, dailyUsage DailyUsageResponse) (*UserUsage, error) {
	userUsage := UserUsage{
		User: user,
	}
	for _, datum := range dailyUsage.Data {
		nTokens := datum.NContextTokensTotal + datum.NGeneratedTokensTotal
		if datum.SnapshotID == "gpt-4-0314" && datum.Operation == "completion" {
			userUsage.NGpt4PromptTokens += datum.NContextTokensTotal
			userUsage.NGpt4CompletionTokens += datum.NGeneratedTokensTotal
			userUsage.PriceUsd += float32(datum.NContextTokensTotal)*0.03/1000 + float32(datum.NGeneratedTokensTotal)*0.06/1000
		} else if datum.SnapshotID == "gpt-4-0613" && datum.Operation == "completion" {
			userUsage.NGpt4PromptTokens += datum.NContextTokensTotal
			userUsage.NGpt4CompletionTokens += datum.NGeneratedTokensTotal
			userUsage.PriceUsd += float32(datum.NContextTokensTotal)*0.03/1000 + float32(datum.NGeneratedTokensTotal)*0.06/1000
		} else if datum.SnapshotID == "gpt-3.5-turbo-16k-0613" && datum.Operation == "completion" {
			userUsage.NGpt3PromptTokens += datum.NContextTokensTotal
			userUsage.NGpt3CompletionTokens += datum.NGeneratedTokensTotal
			userUsage.PriceUsd += float32(datum.NContextTokensTotal)*0.003/1000 + float32(datum.NGeneratedTokensTotal)*0.004/1000
		} else if datum.SnapshotID == "gpt-3.5-turbo-0613" && datum.Operation == "completion" {
			userUsage.NGpt3PromptTokens += datum.NContextTokensTotal
			userUsage.NGpt3CompletionTokens += datum.NGeneratedTokensTotal
			userUsage.PriceUsd += float32(datum.NContextTokensTotal)*0.0015/1000 + float32(datum.NGeneratedTokensTotal)*0.002/1000
		} else if datum.SnapshotID == "gpt-3.5-turbo-0301" && datum.Operation == "completion" {
			userUsage.NGpt3PromptTokens += datum.NContextTokensTotal
			userUsage.NGpt3CompletionTokens += datum.NGeneratedTokensTotal
			userUsage.PriceUsd += float32(datum.NContextTokensTotal)*0.0015/1000 + float32(datum.NGeneratedTokensTotal)*0.002/1000
		} else if datum.SnapshotID == "text-davinci:003" && datum.Operation == "completion" {
			userUsage.NDavinciTokens += nTokens
			userUsage.PriceUsd += float32(nTokens) * 0.0200 / 1000
		} else if datum.SnapshotID == "code-davinci-edit:001" && datum.Operation == "edit" {
			continue
		} else if datum.SnapshotID == "text-embedding-ada-002-v2" && datum.Operation == "embeddings" {
			userUsage.NAdaEmbeddingTokens += nTokens
			userUsage.PriceUsd += float32(nTokens) * 0.0004 / 1000
		} else {
			fmt.Printf(
				"calculations for model %s and operation %s not supported",
				datum.SnapshotID,
				datum.Operation,
			)
		}
	}

	return &userUsage, nil
}
