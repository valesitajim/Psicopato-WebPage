package handlers

import (
	"html/template"
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/models"
)

type UserHandler struct {
	tmplLogin  *template.Template
	tmplRegister *template.Template
	store *database.UserStore
}

func NewUserHandler(tmplLogin *template.Template, tmplRegister *template.Template, store *database.UserStore) *UserHandler {
	return &UserHandler{
		tmplLogin:  tmplLogin,
		tmplRegister: tmplRegister,
		store: store,
	}
}

//Métodos de visualización
func (h *UserHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
    h.tmplLogin.ExecuteTemplate(w, "login.html", nil)
}

func (h *UserHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
    h.tmplRegister.ExecuteTemplate(w, "register.html", nil)
}

// SubmitForm procesa el registro de nuevos usuarios
func (h *UserHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	password := r.FormValue("password") // Cogemos la contraseña del HTML

	// 1. CIFRAR LA CONTRASEÑA (Como pide tu profesor, coste 12)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// 2. Crear el usuario con el Hash, NUNCA con la contraseña real
	user := models.User{
		Name:         r.FormValue("nombre"),
		Email:        r.FormValue("email"),
		PasswordHash: string(hash), // Guardamos el resumen
	}

	// Si hay error al guardar:
	if err := h.store.AppendUser(user); err != nil {
		// Redirigimos al registro con un aviso de error
		http.Redirect(w, r, "/registro?error=servidor", http.StatusSeeOther)
		return
	}
	// SI TODO VA BIEN: Redirigimos al login con mensaje de éxito
	http.Redirect(w, r, "/login?exito=registro", http.StatusSeeOther)

}

// Login procesa el intento de entrada
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")

	// 1. Buscar al usuario en la base de datos
	// Si el usuario no existe o la contraseña está mal:
	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/login?error=credenciales", http.StatusSeeOther)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login?error=credenciales", http.StatusSeeOther)
		return
	}

	// SI EL LOGIN ES CORRECTO: Redirigimos a la página principal
	http.Redirect(w, r, "/?exito=login", http.StatusSeeOther)
}