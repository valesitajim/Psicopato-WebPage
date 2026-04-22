package models

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
}