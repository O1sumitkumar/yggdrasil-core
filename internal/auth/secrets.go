package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

// SecretStore persists sensitive values outside SQLite.
type SecretStore struct {
	dir string
}

func NewSecretStore(dataDir string) *SecretStore {
	return &SecretStore{dir: filepath.Join(dataDir, "secrets")}
}

func (s *SecretStore) Write(name, value string) error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return err
	}
	path := filepath.Join(s.dir, name)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(value), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *SecretStore) Read(name string) (string, error) {
	path := filepath.Join(s.dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *SecretStore) Delete(name string) error {
	path := filepath.Join(s.dir, name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *SecretStore) Path(name string) string {
	return filepath.Join(s.dir, name)
}

func (s *SecretStore) EnsureDir() error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create secrets dir: %w", err)
	}
	return nil
}
