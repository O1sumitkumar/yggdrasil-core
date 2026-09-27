package repositories_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/store"
	"github.com/yeixio/yggdrasil-core/internal/store/repositories"
)

func TestConversationDeleteAndModel(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	repo := repositories.NewConversationRepo(db.SQL)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Hello", "general-assistant", "model-a")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ModelID != "model-a" {
		t.Fatalf("model_id = %q", created.ModelID)
	}

	_, err = repo.AddMessage(ctx, created.ID, "user", "hi")
	if err != nil {
		t.Fatalf("add message: %v", err)
	}

	modelB := "model-b"
	updated, err := repo.Update(ctx, created.ID, repositories.ConversationPatch{ModelID: &modelB})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.ModelID != "model-b" {
		t.Fatalf("updated model_id = %q", updated.ModelID)
	}
	if updated.ProfileID != "general-assistant" {
		t.Fatalf("profile should be unchanged, got %q", updated.ProfileID)
	}

	msgs, err := repo.ListMessages(ctx, created.ID)
	if err != nil || len(msgs) != 1 {
		t.Fatalf("messages before delete: %v len=%d", err, len(msgs))
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	msgs, err = repo.ListMessages(ctx, created.ID)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(msgs) != 0 {
		t.Fatalf("expected cascade delete of messages, got %d", len(msgs))
	}
	if _, err := repo.Get(ctx, created.ID); err == nil {
		t.Fatal("expected get after delete to fail")
	}
}
