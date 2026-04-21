package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/models"
)

type UserHandler struct {
	tmpl  *template.Template
	store *database.UserStore
}

func NewUserHandler(tmpl *template.Template, store *database.UserStore) *UserHandler {
	return &UserHandler{
		tmpl:  tmpl,
		store: store,
	}
}

// ShowForm muestra el formulario de login/registro
func (h *UserHandler) ShowForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := h.tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error renderizando plantilla", http.StatusInternalServerError)
	}
}

// SubmitForm procesa el registro de nuevos usuarios
func (h *UserHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:  r.FormValue("name"), // Asegúrate de que coincida con name="name" en tu HTML
		Email: r.FormValue("email"),
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Nombre y email son obligatorios", http.StatusBadRequest)
		return
	}

	// Usamos AppendUser (o SaveUser según lo tengas en tu database)
	if err := h.store.AppendUser(user); err != nil {
		http.Error(w, "No se pudo guardar el usuario", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "¡Usuario %s guardado correctamente!", user.Name)
}

// Login procesa el intento de entrada
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// USAMOS las variables para que Go no dé error
	log.Printf("Intento de login: Email=%s, Password=%s", email, password)
	
	fmt.Fprintf(w, "¡Bienvenido de nuevo! Has iniciado sesión con: %s", email)
}