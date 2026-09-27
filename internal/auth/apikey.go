package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const apiKeyPrefix = "ygg_"

// APIKeyManager handles API key lifecycle.
type APIKeyManager struct {
	db      *sql.DB
	secrets *SecretStore
}

func NewAPIKeyManager(db *sql.DB, secrets *SecretStore) *APIKeyManager {
	return &APIKeyManager{db: db, secrets: secrets}
}

// APIKeyRecord is the public metadata for a key.
type APIKeyRecord struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsed  *time.Time `json:"last_used_at,omitempty"`
	Revoked   bool       `json:"revoked"`
}

// Create generates a new API key; secret returned once.
func (m *APIKeyManager) Create(ctx context.Context, name string) (record APIKeyRecord, secret string, err error) {
	if name == "" {
		name = "default"
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return record, "", err
	}
	secret = apiKeyPrefix + hex.EncodeToString(raw)
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return record, "", err
	}
	id := uuid.NewString()
	prefix := secret[:12]
	now := time.Now().UTC()
	_, err = m.db.ExecContext(ctx, `
		INSERT INTO api_keys (id, name, key_prefix, key_hash, created_at)
		VALUES (?, ?, ?, ?, ?)`, id, name, prefix, string(hash), now.Format(time.RFC3339Nano))
	if err != nil {
		return record, "", err
	}
	if err := m.secrets.Write("apikey-"+id, secret); err != nil {
		return record, "", err
	}
	return APIKeyRecord{ID: id, Name: name, Prefix: prefix, CreatedAt: now}, secret, nil
}

// Verify checks a presented API key.
func (m *APIKeyManager) Verify(ctx context.Context, secret string) (APIKeyRecord, error) {
	if len(secret) < 12 {
		return APIKeyRecord{}, fmt.Errorf("invalid api key")
	}
	prefix := secret[:12]
	row := m.db.QueryRowContext(ctx, `
		SELECT id, name, key_prefix, key_hash, created_at, revoked_at, last_used_at
		FROM api_keys WHERE key_prefix = ? AND revoked_at IS NULL`, prefix)
	var rec APIKeyRecord
	var hash, created string
	var revoked, lastUsed sql.NullString
	if err := row.Scan(&rec.ID, &rec.Name, &rec.Prefix, &hash, &created, &revoked, &lastUsed); err != nil {
		return APIKeyRecord{}, fmt.Errorf("invalid api key")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret)); err != nil {
		return APIKeyRecord{}, fmt.Errorf("invalid api key")
	}
	rec.CreatedAt = parseTime(created)
	if lastUsed.Valid {
		t := parseTime(lastUsed.String)
		rec.LastUsed = &t
	}
	_, _ = m.db.ExecContext(ctx, `UPDATE api_keys SET last_used_at = ? WHERE id = ?`, time.Now().UTC().Format(time.RFC3339Nano), rec.ID)
	return rec, nil
}

// List returns non-revoked keys (metadata only).
func (m *APIKeyManager) List(ctx context.Context) ([]APIKeyRecord, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT id, name, key_prefix, created_at, revoked_at, last_used_at
		FROM api_keys ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []APIKeyRecord
	for rows.Next() {
		var rec APIKeyRecord
		var created string
		var revoked, lastUsed sql.NullString
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Prefix, &created, &revoked, &lastUsed); err != nil {
			return nil, err
		}
		rec.CreatedAt = parseTime(created)
		rec.Revoked = revoked.Valid
		if lastUsed.Valid {
			t := parseTime(lastUsed.String)
			rec.LastUsed = &t
		}
		if !rec.Revoked {
			out = append(out, rec)
		}
	}
	if out == nil {
		out = []APIKeyRecord{}
	}
	return out, rows.Err()
}

// Revoke invalidates a key.
func (m *APIKeyManager) Revoke(ctx context.Context, id string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := m.db.ExecContext(ctx, `UPDATE api_keys SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`, now, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("api key %q not found", id)
	}
	_ = m.secrets.Delete("apikey-" + id)
	return nil
}

// Rotate revokes old key and creates a new one with same name.
func (m *APIKeyManager) Rotate(ctx context.Context, id string) (APIKeyRecord, string, error) {
	var name string
	err := m.db.QueryRowContext(ctx, `SELECT name FROM api_keys WHERE id = ?`, id).Scan(&name)
	if err != nil {
		return APIKeyRecord{}, "", fmt.Errorf("api key %q not found", id)
	}
	if err := m.Revoke(ctx, id); err != nil {
		return APIKeyRecord{}, "", err
	}
	return m.Create(ctx, name)
}

// Fingerprint returns sha256 hex of secret for logging.
func Fingerprint(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func parseTime(s string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
