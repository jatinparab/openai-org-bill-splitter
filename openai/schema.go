package openai

type User struct {
	Object  string `json:"object"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Member struct {
	Object string `json:"object"`
	Data   []struct {
		Object  string `json:"object"`
		Role    string `json:"role"`
		Created int    `json:"created"`
		User    User   `json:"user"`
	} `json:"data"`
}

type UsersResponse struct {
	Members   Member        `json:"members"`
	Invited   []interface{} `json:"invited"`
	CanInvite bool          `json:"can_invite"`
}

type DailyUsageResponse struct {
	Object string            `json:"object"`
	Data   []DailyUsageDatum `json:"data"`
}

type DailyUsageDatum struct {
	AggregationTimestamp  int    `json:"aggregation_timestamp"`
	NRequests             int    `json:"n_requests"`
	Operation             string `json:"operation"`
	SnapshotID            string `json:"snapshot_id"`
	NContext              int    `json:"n_context"`
	NContextTokensTotal   int    `json:"n_context_tokens_total"`
	NGenerated            int    `json:"n_generated"`
	NGeneratedTokensTotal int    `json:"n_generated_tokens_total"`
}

type UserUsage struct {
	User                  User
	Date                  string
	NGpt4PromptTokens     int
	NGpt4CompletionTokens int
	NGpt3PromptTokens     int
	NGpt3CompletionTokens int
	NDavinciTokens        int
	NAdaEmbeddingTokens   int
	PriceUsd              float32
}
