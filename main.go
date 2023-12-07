package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jatinparab98/openai-org-bill-splitter/lib"
	"github.com/jatinparab98/openai-org-bill-splitter/openai"
	"github.com/joho/godotenv"
)

func main() {
	// Load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Get secrets from env vars
	orgId := os.Getenv("OPENAI_ORG_ID")
	usersResponse, err := openai.GetOrgUsers(orgId)
	if err != nil {
		log.Fatal(err)
	}
	dates := lib.GetDates(time.September, 2023)
	userMap := make(map[openai.User]float32)
	reqCount := 0
	for _, date := range dates {
		userUsages := []openai.UserUsage{}
		for _, v := range usersResponse.Members.Data {
			usage, err := openai.GetDayUsage(v.User, date)
			reqCount += 1
			if reqCount%5 == 0 {
				time.Sleep(time.Minute)
			}
			if err != nil {

				log.Fatal(err)
			}
			userUsage, err := openai.CalculateUserUsage(v.User, *usage)
			if err != nil {
				log.Fatal(err)
			}
			if userUsage.PriceUsd > 0 {
				userUsages = append(userUsages, *userUsage)
			}
			userMap[v.User] += userUsage.PriceUsd
		}
		if len(userUsages) > 0 {
			fmt.Println(date, "------")
			for _, userUsage := range userUsages {
				fmt.Printf("Name: %s\n", userUsage.User.Name)
				if userUsage.NGpt3CompletionTokens > 0 {
					fmt.Printf("GPT 3 Completion Tokens: %v\n", userUsage.NGpt3CompletionTokens)
				}
				if userUsage.NGpt3PromptTokens > 0 {
					fmt.Printf("GPT 3 Prompt Tokens: %v\n", userUsage.NGpt3PromptTokens)
				}
				if userUsage.NGpt4CompletionTokens > 0 {
					fmt.Printf("GPT 4 Completion Tokens: %v\n", userUsage.NGpt4CompletionTokens)
				}
				if userUsage.NGpt4PromptTokens > 0 {
					fmt.Printf("GPT 4 Prompt Tokens: %v\n", userUsage.NGpt4PromptTokens)
				}
				if userUsage.NDavinciTokens > 0 {
					fmt.Printf("Da vinci tokens: %v\n", userUsage.NDavinciTokens)
				}
				if userUsage.NAdaEmbeddingTokens > 0 {
					fmt.Printf("Ada embeddings tokens: %v\n", userUsage.NAdaEmbeddingTokens)
				}
				fmt.Printf("Price USD: %.2f\n", userUsage.PriceUsd)
				fmt.Println()
			}
			fmt.Println()
			fmt.Println("-------")
		}
	}

	// Final total bill calculations
	fmt.Println("Totals")
	fmt.Println("----")
	usdToInr := 82.62
	fmt.Printf("USD to INR Rate: %.2f\n", usdToInr)
	fmt.Println("----")
	var totalOrgUsdPrice float32
	nUsersWithBill := 0
	for _, totalPriceUsd := range userMap {
		totalOrgUsdPrice += totalPriceUsd
		if totalPriceUsd > 0 {
			nUsersWithBill += 1
		}
	}
	fmt.Printf("Total org bill USD Without Tax ->> %.2f\n", totalOrgUsdPrice)

	// 18 % GST
	tax := totalOrgUsdPrice * 18 / 100
	fmt.Printf("Tax ->> %.2f\n", tax)
	var totalOrgUsdPriceWithTax float32
	for user, totalPriceUsdWithoutTax := range userMap {
		if !(totalPriceUsdWithoutTax > 0) {
			continue
		}
		totalPriceUsd := totalPriceUsdWithoutTax + (tax * float32(totalPriceUsdWithoutTax/totalOrgUsdPrice))
		totalOrgUsdPriceWithTax += totalPriceUsd
		totalPriceInr := totalPriceUsd * float32(usdToInr)
		fmt.Printf("%s: %.2f USD ->> %.2f INR\n", user.Name, totalPriceUsd, totalPriceInr)
	}
	fmt.Printf("Total org bill USD ->> %.2f\n", totalOrgUsdPriceWithTax)
}
