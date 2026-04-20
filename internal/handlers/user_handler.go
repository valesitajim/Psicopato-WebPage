package handlers //procesan las peticiones

import (
	"html/template"
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

type formViewData struct {
	Message string
}

// GET /
func (h *UserHandler) ShowForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	data := formViewData{
		Message: r.URL.Query().Get("msg"),
	}

	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error renderizando plantilla", http.StatusInternalServerError)
	}
}

// POST /submit
func (h *UserHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Nombre y email son obligatorios", http.StatusBadRequest)
		return
	}

	if err := h.store.AppendUser(user); err != nil {
		http.Error(w, "No se pudo guardar el usuario", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?msg=Usuario+guardado+correctamente", http.StatusSeeOther)
}
