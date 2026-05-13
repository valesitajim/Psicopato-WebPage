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

// GetUserByEmail busca un  usuarioen el archivo JSONL
func (s *UserStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Abrimos el archivo para leerlo
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf(" Usuario no encontrado")
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
	return nil, fmt.Errorf(" Usuario no encontrado") // Terminamos de leer y no estaba
}

// 1. OBTENER TODOS LOS USUARIOS
func (s *UserStore) GetAllUsers() ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var usuarios []models.User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var u models.User
		if err := json.Unmarshal(scanner.Bytes(), &u); err == nil {
			usuarios = append(usuarios, u)
		}
	}
	return usuarios, nil
}

// 2. OBTENER UN USUARIO POR ID
func (s *UserStore) GetUserByID(id int) (models.User, error) {
	usuarios, err := s.GetAllUsers()
	if err != nil {
		return models.User{}, err
	}

	for _, u := range usuarios {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, fmt.Errorf("usuario no encontrado")
}

// 3. AÑADIR UN NUEVO USUARIO (con ID Autoincremental)
func (s *UserStore) AddUser(nuevoUsuario models.User) (models.User, error) {
	//Bloqueamos el Mutex ANTES de leer y escribir. 
	// Así nadie más puede tocar el archivo mientras calculamos el nuevo ID.
	s.mu.Lock()
	defer s.mu.Unlock()

	//Leemos el archivo para buscar cuál es el ID más alto actualmente
	maxID := 0
	fileLectura, err := os.Open(s.filePath)
	if err == nil {
		scanner := bufio.NewScanner(fileLectura)
		for scanner.Scan() {
			var u models.User
			if err := json.Unmarshal(scanner.Bytes(), &u); err == nil {
				if u.ID > maxID {
					maxID = u.ID // Actualizamos el máximo encontrado
				}
			}
		}
		fileLectura.Close()
	}

	//Asignamos el nuevo ID real correlativo (el máximo + 1)
	nuevoUsuario.ID = maxID + 1

	//Abrimos el archivo en modo "Añadir" (Append) para guardar al nuevo usuario
	fileEscritura, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return models.User{}, err
	}
	defer fileEscritura.Close()

	//Convertimos a JSON y guardamos
	datosJSON, _ := json.Marshal(nuevoUsuario)
	_, err = fileEscritura.WriteString(string(datosJSON) + "\n")
	
	return nuevoUsuario, err
}

//ACTUALIZAR UN USUARIO
func (s *UserStore) UpdateUser(id int, datosActualizados models.User) (models.User, error) {
	usuarios, err := s.GetAllUsers() // Leemos todos sin bloqueo porque GetAllUsers ya lo tiene
	if err != nil {
		return models.User{}, err
	}

	encontrado := false
	for i, u := range usuarios {
		if u.ID == id {
			datosActualizados.ID = id // Mantenemos el ID original
			usuarios[i] = datosActualizados
			encontrado = true
			break
		}
	}

	if !encontrado {
		return models.User{}, fmt.Errorf("usuario no encontrado")
	}

	// Reescribimos el archivo completo
	return datosActualizados, s.reescribirArchivo(usuarios)
}

// 5. BORRAR UN USUARIO
func (s *UserStore) DeleteUser(id int) error {
	usuarios, err := s.GetAllUsers()
	if err != nil {
		return err
	}

	var usuariosRestantes []models.User
	encontrado := false
	for _, u := range usuarios {
		if u.ID != id {
			usuariosRestantes = append(usuariosRestantes, u)
		} else {
			encontrado = true
		}
	}

	if !encontrado {
		return fmt.Errorf("usuario no encontrado")
	}

	return s.reescribirArchivo(usuariosRestantes)
}

// FUNCIÓN DE APOYO: Reescribe el archivo JSONL entero (necesaria para Update y Delete)
func (s *UserStore) reescribirArchivo(usuarios []models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// O_TRUNC vacía el archivo antes de escribir
	file, err := os.OpenFile(s.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, u := range usuarios {
		datosJSON, _ := json.Marshal(u)
		file.WriteString(string(datosJSON) + "\n")
	}
	return nil
}