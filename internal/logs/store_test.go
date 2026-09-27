package logs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListAndReadTail(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "daemon.log"), []byte("line1\nline2\nline3\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "llamacpp-abc.log"), []byte("runtime ok\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "secret.key"), []byte("nope"), 0o600)

	s := &Store{Dir: dir}
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(list))
	}

	content, err := s.ReadTail("daemon.log", 1024)
	if err != nil {
		t.Fatal(err)
	}
	if content.Content == "" || content.Kind != "daemon" {
		t.Fatalf("unexpected content: %#v", content)
	}

	if _, err := s.ReadTail("../escape.log", 1024); err == nil {
		t.Fatal("expected path traversal rejection")
	}
	if _, err := s.ReadTail("secret.key", 1024); err == nil {
		t.Fatal("expected secret rejection")
	}
}

func TestReadTailTruncates(t *testing.T) {
	dir := t.TempDir()
	body := "alpha\n" + strings.Repeat("x", 200) + "\nomega\n"
	_ = os.WriteFile(filepath.Join(dir, "daemon.log"), []byte(body), 0o644)
	s := &Store{Dir: dir}
	content, err := s.ReadTail("daemon.log", 50)
	if err != nil {
		t.Fatal(err)
	}
	if !content.Truncated {
		t.Fatal("expected truncated")
	}
	if !strings.Contains(content.Content, "omega") {
		t.Fatalf("expected tail content, got %q", content.Content)
	}
}
