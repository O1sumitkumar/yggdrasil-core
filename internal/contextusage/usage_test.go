package contextusage

import (
	"testing"

	"github.com/yeixio/yggdrasil-core/pkg/contracts"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

func TestMeasureScalesToRuntimeTokens(t *testing.T) {
	usage := Measure("instructions", "tools", []pluginapi.ChatMessage{
		{Role: "system", Content: "ignored because the pieces are passed separately"},
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "call"},
		{Role: "user", Content: "Tool result for you, not for the user:\npage"},
	}, 100)
	sum := usage.Instructions + usage.Tools + usage.Conversation + usage.ToolResults
	if sum != 100 || usage.PromptTokens != 100 || usage.Estimated {
		t.Fatalf("usage=%+v sum=%d", usage, sum)
	}
	if usage.ToolResults == 0 || usage.Instructions == 0 || usage.Tools == 0 || usage.Conversation == 0 {
		t.Fatalf("missing section: %+v", usage)
	}
}

func TestMeasureEstimatesWithoutRuntimeTokens(t *testing.T) {
	usage := Measure("1234", "", []pluginapi.ChatMessage{{Role: "user", Content: "12345678"}}, 0)
	if !usage.Estimated || usage.Instructions != 1 || usage.Conversation != 2 || usage.PromptTokens != 3 {
		t.Fatalf("%+v", usage)
	}
}

func TestWithoutCurrentTurnKeepsEarlierCopy(t *testing.T) {
	stored := []contracts.Message{
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "hi"},
		{Role: "user", Content: "hello"},
	}
	got := WithoutCurrentTurn(stored, "hello")
	if len(got) != 2 || got[0].Content != "hello" || got[1].Content != "hi" {
		t.Fatalf("%+v", got)
	}
}

func TestFitPriorDropsOldest(t *testing.T) {
	prior := []pluginapi.ChatMessage{
		{Role: "user", Content: "aaaa"},
		{Role: "assistant", Content: "bbbb"},
		{Role: "user", Content: "cccc"},
	}
	got := FitPrior(prior, 8)
	if len(got) != 2 || got[0].Content != "bbbb" || got[1].Content != "cccc" {
		t.Fatalf("%+v", got)
	}
	if FitPrior(prior, 0) != nil {
		t.Fatal("no room should keep nothing")
	}
}
