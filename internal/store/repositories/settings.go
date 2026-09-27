package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

// SettingsRepo stores key/value settings in SQLite.
type SettingsRepo struct {
	db *sql.DB
}

func NewSettingsRepo(db *sql.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

func (r *SettingsRepo) Get(ctx context.Context, key string) (string, bool, error) {
	var value string
	err := r.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (r *SettingsRepo) Set(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}

func (r *SettingsRepo) GetBool(ctx context.Context, key string, def bool) (bool, error) {
	v, ok, err := r.Get(ctx, key)
	if err != nil || !ok {
		return def, err
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def, nil
	}
	return b, nil
}

func (r *SettingsRepo) SetBool(ctx context.Context, key string, value bool) error {
	return r.Set(ctx, key, strconv.FormatBool(value))
}

func (r *SettingsRepo) GetString(ctx context.Context, key, def string) (string, error) {
	v, ok, err := r.Get(ctx, key)
	if err != nil || !ok || v == "" {
		return def, err
	}
	return v, nil
}

func (r *SettingsRepo) GetInt(ctx context.Context, key string, def int) (int, error) {
	v, ok, err := r.Get(ctx, key)
	if err != nil || !ok {
		return def, err
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def, nil
	}
	return n, nil
}

func (r *SettingsRepo) SetInt(ctx context.Context, key string, value int) error {
	return r.Set(ctx, key, strconv.Itoa(value))
}
