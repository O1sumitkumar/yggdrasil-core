package terminal

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteRequiresCommand(t *testing.T) {
	if _, err := New().Execute(context.Background(), map[string]any{"command": "  "}); err == nil {
		t.Fatal("blank command was accepted")
	}
}

func TestExecuteReportsOutputAndExit(t *testing.T) {
	tool := New()
	ok, err := tool.Execute(context.Background(), map[string]any{"command": "printf hi"})
	if err != nil {
		t.Fatal(err)
	}
	if ok["stdout"] != "hi" || ok["exit"] != 0 {
		t.Fatalf("result=%v", ok)
	}

	failed, err := tool.Execute(context.Background(), map[string]any{"command": "exit 3"})
	if err != nil {
		t.Fatal(err)
	}
	if failed["exit"] != 3 {
		t.Fatalf("result=%v", failed)
	}
	if !strings.Contains(failed["error"].(string), "exit status 3") {
		t.Fatalf("error=%v", failed["error"])
	}
}
