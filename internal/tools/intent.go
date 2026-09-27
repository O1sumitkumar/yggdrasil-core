package tools

import "strings"

// MessageNeedsLiveWeb reports questions that need current public information.
func MessageNeedsLiveWeb(message string) bool {
	value := strings.ToLower(message)
	needles := []string{
		"weather", "forecast", "latest version", "latest news", "who won",
		"last night", "right now", "today", "look up", "documentation",
		"summarize this url", "summarize the url", "http://", "https://",
		"current price", "stock price", "search the web", "search online",
	}
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}
