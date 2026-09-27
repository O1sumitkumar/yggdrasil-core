package tools

import (
	"strings"
	"testing"
)

func TestStructuredToolCallIsNotVisible(t *testing.T) {
	out := ParseModelOutput(`{"tool_call":{"id":"internet.search","args":{"query":"Juneau weather"}}}`)
	if out.Call == nil || out.Call.ID != "internet.search" || out.Call.Args["query"] != "Juneau weather" {
		t.Fatalf("call=%+v", out.Call)
	}
	if strings.Contains(out.Text, "tool_call") || strings.Contains(out.Text, "Juneau") {
		t.Fatalf("visible %q", out.Text)
	}
}

func TestTextEncodedToolCallKeepsProse(t *testing.T) {
	out := ParseModelOutput("Let me check that.\n" + `{"tool_call":{"id":"internet.search","args":{"query":"Juneau weather"}}}`)
	if out.Call == nil || out.Call.ID != "internet.search" {
		t.Fatalf("call=%+v", out.Call)
	}
	if out.Text != "Let me check that." || strings.Contains(out.Text, "tool_call") {
		t.Fatalf("visible %q", out.Text)
	}
}

func TestOrdinaryJSONStaysVisible(t *testing.T) {
	sample := "Here is the object:\n{\"name\":\"Ada\",\"year\":1815}"
	out := ParseModelOutput(sample)
	if out.Call != nil {
		t.Fatalf("ordinary JSON was a tool call: %+v", out.Call)
	}
	if !strings.Contains(out.Text, `"name"`) || !strings.Contains(out.Text, "1815") {
		t.Fatalf("visible %q", out.Text)
	}
}

func TestCodeFenceJSONStaysVisible(t *testing.T) {
	sample := "Example:\n```json\n{\"query\":\"not a tool\"}\n```"
	out := ParseModelOutput(sample)
	if out.Call != nil || !strings.Contains(out.Text, `"query"`) {
		t.Fatalf("out=%+v text=%q", out.Call, out.Text)
	}
}

func TestMalformedToolCallDoesNotExecute(t *testing.T) {
	out := ParseModelOutput(`{"tool_call":{"id":"internet.search","args":`)
	if out.Call != nil {
		t.Fatalf("executed partial call %+v", out.Call)
	}
	if !out.Malformed || strings.Contains(out.Text, "tool_call") {
		t.Fatalf("malformed=%v text=%q", out.Malformed, out.Text)
	}
}

func TestUnknownToolAndPartialArgsAreRejected(t *testing.T) {
	unknown := ParseModelOutput(`{"tool_call":{"id":"weather_check","args":{"city":"Juneau"}}}`)
	if unknown.Call != nil || strings.Contains(unknown.Text, "weather_check") {
		t.Fatalf("unknown call leaked: %+v %q", unknown.Call, unknown.Text)
	}
	missing := ParseModelOutput(`{"tool_call":{"id":"internet.search","args":{}}}`)
	if missing.Call != nil {
		t.Fatal("missing query was executed")
	}
}

func TestToolFailureEnvelopeIsHidden(t *testing.T) {
	raw := "Juneau is currently cloudy.\n" + `{"tool_result":{"tool_id":"internet.search","ok":false,"error":"stack trace line"}}`
	out := ParseModelOutput(raw)
	if strings.Contains(out.Text, "tool_result") || strings.Contains(out.Text, "stack trace") {
		t.Fatalf("visible %q", out.Text)
	}
	if !strings.Contains(out.Text, "cloudy") {
		t.Fatalf("visible %q", out.Text)
	}
}

func TestToolNarrationIsRemovedWhenACallIsPresent(t *testing.T) {
	raw := `I will now use the web tool. {"tool_call":{"id":"internet.search","args":{"query":"Juneau weather"}}}`
	out := ParseModelOutput(raw)
	if out.Call == nil {
		t.Fatal("expected a call")
	}
	if strings.Contains(strings.ToLower(out.Text), "web tool") || strings.Contains(out.Text, "tool_call") {
		t.Fatalf("visible %q", out.Text)
	}
}

func TestTaggedAndFunctionCallFormats(t *testing.T) {
	tagged := ParseModelOutput(`<tool_call>{"id":"internet.open","args":{"url":"https://weather.gov"}}</tool_call>`)
	if tagged.Call == nil || tagged.Call.ID != "internet.open" || tagged.Text != "" {
		t.Fatalf("tagged %+v %q", tagged.Call, tagged.Text)
	}
	fn := ParseModelOutput(`function_call(git.status, {})`)
	if fn.Call == nil || fn.Call.ID != "git.status" || strings.Contains(fn.Text, "function_call") {
		t.Fatalf("function %+v %q", fn.Call, fn.Text)
	}
}
