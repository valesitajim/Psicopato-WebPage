package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"TIENDAPATOS/internal/models"
)

type UserStore struct {
	filePath string
	mu       sync.Mutex
}

func NewUserStore(filePath string) *UserStore {
	return &UserStore{filePath: filePath}
}

func (s *UserStore) AppendUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(user) // Cada Encode añade una línea JSON (JSONL)
}
