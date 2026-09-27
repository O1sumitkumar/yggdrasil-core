package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NodeIdentity holds ed25519 node keys.
type NodeIdentity struct {
	NodeID     string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	CertPEM    []byte
}

// LoadOrCreateIdentity loads or generates node identity keys.
func LoadOrCreateIdentity(secrets *SecretStore, nodeID string) (*NodeIdentity, error) {
	if err := secrets.EnsureDir(); err != nil {
		return nil, err
	}
	privPath := secrets.Path("node-" + nodeID + ".key")
	pubPath := secrets.Path("node-" + nodeID + ".pub")

	if data, err := os.ReadFile(privPath); err == nil {
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, fmt.Errorf("invalid private key pem")
		}
		priv := ed25519.PrivateKey(block.Bytes)
		pubData, err := os.ReadFile(pubPath)
		if err != nil {
			return nil, err
		}
		pubBlock, _ := pem.Decode(pubData)
		if pubBlock == nil {
			return nil, fmt.Errorf("invalid public key pem")
		}
		return &NodeIdentity{
			NodeID:     nodeID,
			PrivateKey: priv,
			PublicKey:  ed25519.PublicKey(pubBlock.Bytes),
			CertPEM:    pubData,
		}, nil
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "ED25519 PRIVATE KEY", Bytes: priv})
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "ED25519 PUBLIC KEY", Bytes: pub})
	if err := os.WriteFile(privPath, privPEM, 0o600); err != nil {
		return nil, err
	}
	if err := os.WriteFile(pubPath, pubPEM, 0o644); err != nil {
		return nil, err
	}
	return &NodeIdentity{NodeID: nodeID, PrivateKey: priv, PublicKey: pub, CertPEM: pubPEM}, nil
}

// Fingerprint returns a stable node certificate fingerprint.
func (id *NodeIdentity) Fingerprint() string {
	return Fingerprint(string(id.PublicKey))
}

// Sign signs message with node private key.
func (id *NodeIdentity) Sign(msg []byte) []byte {
	return ed25519.Sign(id.PrivateKey, msg)
}

// VerifyPeer verifies a peer signature with their public key PEM.
func VerifyPeer(pubPEM, msg, sig []byte) bool {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// ed25519 raw key
		if len(block.Bytes) == ed25519.PublicKeySize {
			return ed25519.Verify(ed25519.PublicKey(block.Bytes), msg, sig)
		}
		return false
	}
	if pk, ok := pub.(ed25519.PublicKey); ok {
		return ed25519.Verify(pk, msg, sig)
	}
	return false
}

// TrustRecord stored in node_trust table.
type TrustRecord struct {
	NodeID      string
	Fingerprint string
	CertPEM     string
	PairedAt    time.Time
}

// WriteTrustFile stores peer cert for debugging (optional).
func WriteTrustFile(dataDir, nodeID string, certPEM []byte) error {
	dir := filepath.Join(dataDir, "trust")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, nodeID+".pem"), certPEM, 0o644)
}

// NowUTC returns current UTC time.
func NowUTC() time.Time { return time.Now().UTC() }
