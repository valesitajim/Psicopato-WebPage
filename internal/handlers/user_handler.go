package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/models"
)

type UserHandler struct {
	templates *template.Template // Cargaremos todos los archivos aquí
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
	// Usamos la lógica de renderizado que ya tenéis en el proyecto
	// Asumo que vuestro sistema carga primero el layout y luego la página específica
	err := h.render(w, r, "admin.html", nil) 
	if err != nil {
		http.Error(w, "Error al cargar la página de administración", http.StatusInternalServerError)
	}
}

// Métodos de visualización
func (h *UserHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "login.html", nil)
}

// Método para mostrar el perfil dinámico
func (h *UserHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	// 1. Leemos la cookie para saber quién es
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Redirect(w, r, "/login?error=no_autorizado", http.StatusSeeOther)
		return
	}

	// 2. Buscamos al User en la base de datos (tu JSON)
	email := cookie.Value
	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/login?error=no_autorizado", http.StatusSeeOther)
		return
	}

	// 3. Empaquetamos los datos dinámicos
	data := map[string]interface{}{
		"Nombre": user.Name,
		"Email":  user.Email,
	}

	// 4. Inyectamos los datos en la plantilla perfil.html
	// Fíjate que aquí NO pasamos nil, pasamos 'data'
	h.templates.ExecuteTemplate(w, "perfil.html", data)
}

func (h *UserHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "register.html", nil)
}

// SubmitForm procesa el registro de nuevos Users
func (h *UserHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	password := r.FormValue("password") // Cogemos la contraseña del HTML

	//CIFRAR LA CONTRASEÑA
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	//Crear el User con el Hash
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

	// 1. Buscar al User en la base de datos
	// Si el User no existe o la contraseña está mal:
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

	// Creamos la cookie de sesión
	cookie := &http.Cookie{
		Name:     "session_user",
		Value:    email, // Guardamos el email como identificador
		Path:     "/",
		HttpOnly: true, // Seguridad: impide que JavaScript acceda a la cookie
		MaxAge:   3600, // Dura 1 hora
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	log.Printf("INFO: Cookie de sesión creada para %s", email)

	http.Redirect(w, r, "/?exito=login", http.StatusSeeOther)

}

// Middleware para proteger rutas
func (h *UserHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_user")
		if err != nil {
			// Si no hay cookie, lo mandamos al login
			log.Printf("WARN: Intento de acceso no autorizado a %s", r.URL.Path)
			http.Redirect(w, r, "/login?error=no_autorizado", http.StatusSeeOther)
			return
		}

		// Si hay cookie, dejamos que pase a la siguiente función (next)
		log.Printf("INFO: User %s accediendo a ruta protegida", cookie.Value)
		next.ServeHTTP(w, r)
	}
}

// Logout cierra la sesión del User eliminando la cookie
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Verificamos que sea un método POST por seguridad
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Sobrescribimos la cookie actual con una caducada
	cookie := &http.Cookie{
		Name:     "session_user",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,                             // Obliga al navegador a borrarla al instante
		Expires:  time.Now().Add(-1 * time.Hour), // Fecha en el pasado por si acaso
	}
	http.SetCookie(w, cookie)

	log.Println("INFO: Un User ha cerrado sesión correctamente.")

	// Redirigimos a la página de login con un mensaje de éxito
	http.Redirect(w, r, "/login?exito=logout", http.StatusSeeOther)
}

func (h *UserHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_user")
	isLoggedIn := (err == nil)

	data := map[string]interface{}{
		"Titulo":    "Inicio - Tienda de Patos",
		"Logueado":  isLoggedIn,
		"UserEmail": "",
	}

	if isLoggedIn {
		data["UserEmail"] = cookie.Value
	}

	//Cargamos el layout Y la página
	// Luego ejecutamos "base", que es el nombre que le dimos al esqueleto
	tmpl, _ := template.ParseFiles("ui/templates/layout.html", "ui/templates/index.html")
	tmpl.ExecuteTemplate(w, "base", data)
}

// GET
func (h *UserHandler) ListarUsers(w http.ResponseWriter, r *http.Request) {
	Users, err := h.store.GetAllUsers() // Lógica real de lectura de tu JSONL
	if err != nil {
		http.Error(w, "Error al leer Users", http.StatusInternalServerError)
		return
	}

	//Establecer la cabecera obligatoria para APIs REST
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	//Codificar a JSON y enviar al navegador
	json.NewEncoder(w).Encode(Users)
}

// GET por ID
func (h *UserHandler) ObtenerUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, _ := strconv.Atoi(idParam)

	User, err := h.store.GetUserByID(id) // Búsqueda real en el archivo
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User no encontrado"})
		return
	}

	//Si existe, establecemos cabeceras y enviamos el JSON con un 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(User)
}

// POST
func (h *UserHandler) CrearUser(w http.ResponseWriter, r *http.Request) {
	//Limitar el tamaño del cuerpo a 1 MB para evitar ataques de denegación de servicio.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var nuevoUser models.User

	//Decodificación estricta del JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Rechaza peticiones con campos extra que no esperas

	if err := decoder.Decode(&nuevoUser); err != nil {
		// Si el JSON está mal formado o tiene campos extra, devolvemos 400 Bad Request
		http.Error(w, "JSON inválido o campos desconocidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	//	Validación de los datos
	if nuevoUser.Name == "" || nuevoUser.Email == "" {
		// Si los datos no son válidos, devolvemos 422 Unprocessable Entity
		http.Error(w, "El nombre y el email son obligatorios", http.StatusUnprocessableEntity)
		return
	}

	nuevoUser, err := h.store.AddUser(nuevoUser) // Guardado real en JSONL
	if err != nil {
		http.Error(w, "Error al guardar", http.StatusInternalServerError)
		return
	}

	// Devolver 201 Created y la cabecera Location
	// Conviene devolver Location con la URL del nuevo recurso
	w.Header().Set("Location", fmt.Sprintf("/api/v1/Users/%d", nuevoUser.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Devolver el objeto completo (con el id generado)
	json.NewEncoder(w).Encode(nuevoUser)
}

// PUT
func (h *UserHandler) ActualizarUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la URL (ej: /api/v1/Users/2)
	idParam := r.PathValue("id")
	id, _ := strconv.Atoi(idParam)

	//Limitar el tamaño del cuerpo a 1 MB para evitar ataques DOS
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var datosActualizados models.User

	//Decodificación estricta del JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Rechaza la petición si tiene campos extra que no están en el struct

	if err := decoder.Decode(&datosActualizados); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest) // 400 Bad Request
		return
	}

	// Validación de datos de negocio
	// La validación en el servidor es obligatoria aunque el frontend también lo valide
	if datosActualizados.Name == "" || datosActualizados.Email == "" {
		http.Error(w, "El nombre y el email son obligatorios", http.StatusUnprocessableEntity) // 422 Unprocessable Entity
		return
	}

	UserActualizado, err := h.store.UpdateUser(id, datosActualizados) // Actualización real
	if err != nil {
		http.Error(w, "User no encontrado", http.StatusNotFound)
		return
	}

	//Devolver 200 OK y el recurso actualizado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserActualizado)
}

// DELETE
func (h *UserHandler) BorrarUser(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el ID de la URL (ej: /api/v1/Users/2)
	idParam := r.PathValue("id")
	id, _ := strconv.Atoi(idParam)

	err := h.store.DeleteUser(id) // Borrado real en el archivo
	if err != nil {
		http.Error(w, "No se pudo eliminar", http.StatusNotFound)
		return
	}
	//Responder éxito sin contenido
	w.WriteHeader(http.StatusNoContent)
}
