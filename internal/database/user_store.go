package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"bufio"
	"fmt"

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

// GetUserByEmail busca un usuario en el archivo JSONL
func (s *UserStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Abrimos el archivo para leerlo
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	defer file.Close()

	// Leemos línea por línea
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var user models.User
		// Convertimos la línea JSON a nuestro struct User
		if err := json.Unmarshal(scanner.Bytes(), &user); err == nil {
			if user.Email == email {
				return &user, nil // ¡Lo encontramos!
			}
		}
	}
	return nil, fmt.Errorf("usuario no encontrado") // Terminamos de leer y no estaba
}