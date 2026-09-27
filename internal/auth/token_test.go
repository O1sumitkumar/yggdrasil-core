package auth

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

func TestAuthTokenRoundTrip(t *testing.T) {
	dir := t.TempDir()
	secrets := NewSecretStore(dir)
	id, err := LoadOrCreateIdentity(secrets, "node-a")
	if err != nil {
		t.Fatal(err)
	}
	token := id.AuthToken()
	nodeID, err := ParseAndVerifyAuthToken(token, id.CertPEM)
	if err != nil {
		t.Fatal(err)
	}
	if nodeID != "node-a" {
		t.Fatalf("nodeID=%q", nodeID)
	}
	peeked, err := PeekAuthTokenNodeID(token)
	if err != nil || peeked != "node-a" {
		t.Fatalf("peek=%q err=%v", peeked, err)
	}
}

func TestAuthTokenRejectsWrongPeer(t *testing.T) {
	dir := t.TempDir()
	a, err := LoadOrCreateIdentity(NewSecretStore(dir+"/a"), "node-a")
	if err != nil {
		t.Fatal(err)
	}
	b, err := LoadOrCreateIdentity(NewSecretStore(dir+"/b"), "node-b")
	if err != nil {
		t.Fatal(err)
	}
	token := a.AuthToken()
	if _, err := ParseAndVerifyAuthToken(token, b.CertPEM); err == nil {
		t.Fatal("expected signature failure against wrong peer cert")
	}
}

func TestAuthTokenRejectsExpired(t *testing.T) {
	dir := t.TempDir()
	id, err := LoadOrCreateIdentity(NewSecretStore(dir), "node-a")
	if err != nil {
		t.Fatal(err)
	}
	exp := time.Now().UTC().Add(-time.Minute).Unix()
	msg := fmt.Sprintf("%s|%d", id.NodeID, exp)
	sig := id.Sign([]byte(msg))
	raw := fmt.Sprintf("%s|%d|%x", id.NodeID, exp, sig)
	token := base64.RawURLEncoding.EncodeToString([]byte(raw))
	if _, err := ParseAndVerifyAuthToken(token, id.CertPEM); err == nil {
		t.Fatal("expected expired token error")
	}
}

func TestAuthTokenRejectsUnsignedGarbage(t *testing.T) {
	if _, err := ParseAndVerifyAuthToken("not-a-token", []byte("x")); err == nil {
		t.Fatal("expected error")
	}
}
