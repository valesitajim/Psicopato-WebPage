package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/models"
)

type UserHandler struct {
	templates *template.Template
	store     *database.UserStore
}

func NewUserHandler(tmpl *template.Template, store *database.UserStore) *UserHandler {
	return &UserHandler{
		templates: tmpl,
		store:     store,
	}
}

// ShowAdmin renderiza la interfaz de administración
func (h *UserHandler) ShowAdmin(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "admin.html", nil)
	if err != nil {
		http.Error(w, "No se pudo cargar la plantilla: "+err.Error(), 500)
	}
}

func (h *UserHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *UserHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.store.GetUserByEmail(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Nombre": user.Name,
		"Email":  user.Email,
	}
	h.templates.ExecuteTemplate(w, "perfil.html", data)
}

func (h *UserHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "register.html", nil)
}

// Registro Normal (Formulario)
func (h *UserHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", 405)
		return
	}
	r.ParseForm()

	hash, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 12)
	user := models.User{
		Name:         r.FormValue("nombre"),
		Email:        r.FormValue("email"),
		PasswordHash: string(hash),
	}

	if err := h.store.AppendUser(user); err != nil {
		http.Redirect(w, r, "/registro?error=1", 303)
		return
	}
	http.Redirect(w, r, "/login?exito=1", 303)
}

// Login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, err := h.store.GetUserByEmail(r.FormValue("email"))
	if err != nil {
		http.Redirect(w, r, "/login?error=1", 303)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.FormValue("password"))); err != nil {
		http.Redirect(w, r, "/login?error=1", 303)
		return
	}

	cookie := &http.Cookie{
		Name: "session_user", Value: user.Email, Path: "/", MaxAge: 3600, HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/?exito=1", 303)
}

func (h *UserHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("session_user"); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{Name: "session_user", Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", 303)
}

func (h *UserHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_user")
	isLoggedIn := (err == nil)

	// Usamos ExecuteTemplate con el nombre del bloque base si lo tienes definido así
	h.templates.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Logueado": isLoggedIn,
	})
}

// --- API REST ---

func (h *UserHandler) ListarUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		http.Error(w, "Error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CrearUser(w http.ResponseWriter, r *http.Request) {
	// 1. Estructura temporal para capturar el JSON del JS
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", 400)
		return
	}

	// 2. Cifrar la contraseña antes de guardar
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

	// 3. Crear el modelo real para el Store
	nuevoUser := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	// 4. Guardar
	userGuardado, err := h.store.AddUser(nuevoUser)
	if err != nil {
		http.Error(w, "Error al guardar", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userGuardado)
}

func (h *UserHandler) BorrarUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, _ := strconv.Atoi(idParam)

	if err := h.store.DeleteUser(id); err != nil {
		http.Error(w, "No encontrado", 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ObtenerUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.store.GetUserByID(id)
	if err != nil {
		http.Error(w, "No encontrado", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ActualizarUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)

	actualizado, err := h.store.UpdateUser(id, u)
	if err != nil {
		http.Error(w, "Error", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actualizado)
}
