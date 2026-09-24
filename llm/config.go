package llm

type Config struct {
	MessagesKey string
	ResponseKey string

	Model       string
	Temperature float64
	MaxTokens   int
}
