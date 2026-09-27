package simple

import (
	"strings"
	"testing"
)

func TestParseToolCall(t *testing.T) {
	id, args, ok := parseToolCall(`{"tool_call":{"id":"filesystem.read","args":{"path":"README.md"}}}`)
	if !ok || id != "filesystem.read" {
		t.Fatalf("got ok=%v id=%q", ok, id)
	}
	if args["path"] != "README.md" {
		t.Fatalf("args=%v", args)
	}

	id, _, ok = parseToolCall("```json\n{\"tool_call\":{\"id\":\"git.status\",\"args\":{}}}\n```")
	if !ok || id != "git.status" {
		t.Fatalf("fenced got ok=%v id=%q", ok, id)
	}

	_, _, ok = parseToolCall("Just a normal answer with no tools.")
	if ok {
		t.Fatal("expected no tool call")
	}
}

func TestParseToolCallMixedWithAnswer(t *testing.T) {
	content := `{"tool_call":{"name":"terminal","arguments":{"command":"echo hi"}}} A typical training plan keeps sessions short.`
	id, args, ok := parseToolCall(content)
	if !ok || id != "terminal" {
		t.Fatalf("got ok=%v id=%q", ok, id)
	}
	if args["command"] != "echo hi" {
		t.Fatalf("args=%v", args)
	}
	if got := userVisibleReply(content); got != "A typical training plan keeps sessions short." {
		t.Fatalf("visible=%q", got)
	}
}

func TestParseToolCallRepairsShellEscapes(t *testing.T) {
	content := `{"tool_call":{"id":"terminal","args":{"command":"grep -Po '(\d+\L*) hour(s?)'"}}} A typical Michael Phelps training plan keeps the work easy.`
	id, args, ok := parseToolCall(content)
	if !ok || id != "terminal" {
		t.Fatalf("got ok=%v id=%q", ok, id)
	}
	command, _ := args["command"].(string)
	if !strings.Contains(command, `\d`) || !strings.Contains(command, `\L`) {
		t.Fatalf("command=%q", command)
	}
	if got := userVisibleReply(content); !strings.HasPrefix(got, "A typical") || strings.Contains(got, "tool_call") {
		t.Fatalf("visible=%q", got)
	}
}

func TestUserVisibleReplyHidesBareToolCall(t *testing.T) {
	content := "```json\n{\"tool_call\":{\"id\":\"filesystem.read\",\"args\":{\"path\":\"README.md\"}}}\n```"
	got := userVisibleReply(content)
	if strings.Contains(got, "tool_call") || strings.Contains(got, "{") {
		t.Fatalf("visible=%q", got)
	}
	if got != "" {
		t.Fatalf("visible=%q", got)
	}
}
