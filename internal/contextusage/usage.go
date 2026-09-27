// Package contextusage splits a prompt into the parts a person can understand:
// instructions, tool definitions, the conversation, and tool results.
package contextusage

import (
	"strings"
	"unicode/utf8"

	"github.com/yeixio/yggdrasil-core/pkg/contracts"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

// ReplyReserveTokens leaves room for the model's answer when older messages are kept.
const ReplyReserveTokens = 512

// DefaultWindow is the token window a local model runs with when its catalog
// window is missing or larger. It matches the llama.cpp start default.
const DefaultWindow = 8192

// Usage is one measured prompt. Section counts sum to PromptTokens.
type Usage struct {
	PromptTokens int  `json:"prompt_tokens"`
	Limit        int  `json:"limit"`
	Instructions int  `json:"instructions"`
	Tools        int  `json:"tools"`
	Conversation int  `json:"conversation"`
	ToolResults  int  `json:"tool_results"`
	Estimated    bool `json:"estimated"`
}

// Map is the chat.complete payload field.
func (u Usage) Map() map[string]any {
	return map[string]any{
		"prompt_tokens": u.PromptTokens,
		"limit":         u.Limit,
		"instructions":  u.Instructions,
		"tools":         u.Tools,
		"conversation":  u.Conversation,
		"tool_results":  u.ToolResults,
		"estimated":     u.Estimated,
	}
}

// Measure sizes a prompt. When the runtime reports prompt tokens, the sections
// are scaled so they add up to that count. Otherwise each section is about one
// token per four characters.
func Measure(instructions, toolPrompt string, messages []pluginapi.ChatMessage, promptTokens int) Usage {
	parts := [4]int{
		utf8.RuneCountInString(instructions),
		utf8.RuneCountInString(toolPrompt),
	}
	for _, msg := range messages {
		if msg.Role == "system" {
			continue
		}
		n := utf8.RuneCountInString(msg.Content)
		if msg.Role == "user" && isToolFeedback(msg.Content) {
			parts[3] += n
			continue
		}
		parts[2] += n
	}
	scaled, total, estimated := scale(parts, promptTokens)
	return Usage{
		PromptTokens: total,
		Instructions: scaled[0],
		Tools:        scaled[1],
		Conversation: scaled[2],
		ToolResults:  scaled[3],
		Estimated:    estimated,
	}
}

func isToolFeedback(content string) bool {
	return strings.HasPrefix(content, "Tool result for you") ||
		strings.HasPrefix(content, "The tool failed.") ||
		strings.HasPrefix(content, "That tool call was not valid")
}

func scale(parts [4]int, promptTokens int) (out [4]int, total int, estimated bool) {
	sum := parts[0] + parts[1] + parts[2] + parts[3]
	if promptTokens <= 0 {
		for i, n := range parts {
			out[i] = n / 4
			total += out[i]
		}
		return out, total, true
	}
	if sum == 0 {
		out[2] = promptTokens
		return out, promptTokens, false
	}
	used := 0
	largest := 0
	for i, n := range parts {
		out[i] = promptTokens * n / sum
		used += out[i]
		if out[i] >= out[largest] {
			largest = i
		}
	}
	out[largest] += promptTokens - used
	return out, promptTokens, false
}

// WithoutCurrentTurn drops the user message just saved for this turn.
func WithoutCurrentTurn(stored []contracts.Message, prompt string) []pluginapi.ChatMessage {
	if len(stored) > 0 {
		last := stored[len(stored)-1]
		if last.Role == "user" && last.Content == prompt {
			stored = stored[:len(stored)-1]
		}
	}
	out := make([]pluginapi.ChatMessage, 0, len(stored))
	for _, msg := range stored {
		if msg.Role != "user" && msg.Role != "assistant" {
			continue
		}
		if strings.TrimSpace(msg.Content) == "" {
			continue
		}
		out = append(out, pluginapi.ChatMessage{Role: msg.Role, Content: msg.Content})
	}
	return out
}

// RoomRunes is how much older chat text can still fit beside the reserved pieces.
func RoomRunes(limitTokens int, reserved ...string) int {
	room := limitTokens*4 - ReplyReserveTokens*4
	for _, text := range reserved {
		room -= utf8.RuneCountInString(text)
	}
	if room < 0 {
		return 0
	}
	return room
}

// FitPrior keeps the newest messages that fit in roomRunes, in chronological order.
func FitPrior(prior []pluginapi.ChatMessage, roomRunes int) []pluginapi.ChatMessage {
	if roomRunes <= 0 || len(prior) == 0 {
		return nil
	}
	start := len(prior)
	used := 0
	for i := len(prior) - 1; i >= 0; i-- {
		n := utf8.RuneCountInString(prior[i].Content)
		if n == 0 {
			continue
		}
		if used+n > roomRunes {
			break
		}
		used += n
		start = i
	}
	if start >= len(prior) {
		return nil
	}
	out := make([]pluginapi.ChatMessage, len(prior)-start)
	copy(out, prior[start:])
	return out
}
