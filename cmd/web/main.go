package main

import (
	"html/template"
	"log"
	"net/http"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/handlers"
)

func main() {

	mux := http.NewServeMux()
	store := database.NewUserStore("api/users.jsonl")

	tmpl, err := template.ParseGlob("ui/templates/*.html")
	if err != nil {
		log.Fatalf("Error cargando templates: %v", err)
	}

	userHandler := handlers.NewUserHandler(tmpl, store)

	// --- RUTAS DE NAVEGACIÓN ---
	mux.HandleFunc("GET /", userHandler.ShowHome)
	mux.HandleFunc("GET /login", userHandler.ShowLogin)
	mux.HandleFunc("GET /registro", userHandler.ShowRegister)
	mux.HandleFunc("GET /perfil", userHandler.AuthMiddleware(userHandler.ShowProfile))

	// Ahora el panel de admin también está protegido
	mux.HandleFunc("GET /admin", userHandler.AuthMiddleware(userHandler.ShowAdmin))

	mux.HandleFunc("GET /chaquetas", func(w http.ResponseWriter, r *http.Request) {
		tmplChaquetas, err := template.ParseFiles("ui/templates/chaquetas_nina.html")
		if err != nil {
			http.Error(w, "Template no encontrado", 404)
			return
		}
		tmplChaquetas.Execute(w, nil)
	})

	// --- PROCESAMIENTO ---
	mux.HandleFunc("POST /procesar-registro", userHandler.SubmitForm)
	mux.HandleFunc("POST /procesar-login", userHandler.Login)
	mux.HandleFunc("GET /logout", userHandler.Logout)

	// --- API (RESTABLECIDA: Con AuthMiddleware para seguridad total) ---
	mux.HandleFunc("GET /api/v1/usuarios", userHandler.AuthMiddleware(userHandler.ListarUsers))
	mux.HandleFunc("POST /api/v1/usuarios", userHandler.AuthMiddleware(userHandler.CrearUser))
	mux.HandleFunc("DELETE /api/v1/usuarios/{id}", userHandler.AuthMiddleware(userHandler.BorrarUser))
	mux.HandleFunc("PUT /api/v1/usuarios/{id}", userHandler.AuthMiddleware(userHandler.ActualizarUser))
	mux.HandleFunc("GET /api/v1/usuarios/{id}", userHandler.AuthMiddleware(userHandler.ObtenerUser))

	// --- ARCHIVOS ESTÁTICOS ---
	fs := http.FileServer(http.Dir("ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	log.Println("Servidor encendido en http://localhost:8080")
	log.Println("Prueba tu panel en http://localhost:8080/admin")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
