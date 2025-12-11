package ports

type LLM interface {
	DoRequest(question string) string
}
