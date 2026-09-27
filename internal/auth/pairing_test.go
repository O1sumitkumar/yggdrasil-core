package auth

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/store"
)

func TestPairingStateMachine(t *testing.T) {
	dir := t.TempDir()
	db, err := store.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	secrets := NewSecretStore(dir)
	id, err := LoadOrCreateIdentity(secrets, "local-node")
	if err != nil {
		t.Fatal(err)
	}
	pm := NewPairingManager(db.SQL, id)

	s, err := pm.StartPairing("remote-1", "Remote PC", "192.168.1.10:7332", []byte("cert-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if s.State != PairingPending {
		t.Fatalf("state=%q", s.State)
	}
	if len(s.Code) != 6 {
		t.Fatalf("code=%q", s.Code)
	}

	approved, err := pm.ApproveByCode(context.Background(), s.Code)
	if err != nil {
		t.Fatal(err)
	}
	if approved.State != PairingApproved {
		t.Fatalf("state=%q", approved.State)
	}

	if err := pm.RevokeTrust(context.Background(), "remote-1"); err != nil {
		t.Fatal(err)
	}
}

func TestPairingRejectInvalidCode(t *testing.T) {
	dir := t.TempDir()
	db, err := store.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	secrets := NewSecretStore(dir)
	id, _ := LoadOrCreateIdentity(secrets, "local")
	pm := NewPairingManager(db.SQL, id)
	if _, err := pm.ApproveByCode(context.Background(), "000000"); err == nil {
		t.Fatal("expected error for invalid code")
	}
}
