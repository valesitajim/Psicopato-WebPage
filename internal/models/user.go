package models

type User struct {
	ID     int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Rol    string `json:"rol,omitempty"` // omitempty lo oculta si está vacío
}