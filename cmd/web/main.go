package main

import (
	"html/template"
	"log"
	"net/http"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/handlers"
)

func main() {
	// 1. Configuración
	store := database.NewUserStore("api/users.jsonl")
	
	// Cargamos el template (ahora compartido para login y registro)
	tmplAuth, err := template.ParseFiles("ui/templates/user_form.html")
	if err != nil {
		log.Fatalf("error cargando template de autenticación: %v", err)
	}
	userHandler := handlers.NewUserHandler(tmplAuth, store)

	// --- RUTAS --- //

	// Inicio
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmplIndex, _ := template.ParseFiles("ui/templates/index.html")
		tmplIndex.Execute(w, nil)
	})

	// Chaquetas
	http.HandleFunc("/chaquetas", func(w http.ResponseWriter, r *http.Request) {
		tmplChaquetas, _ := template.ParseFiles("ui/templates/chaquetas_nina.html")
		tmplChaquetas.Execute(w, nil)
	})

	// RUTAS DE USUARIO
	http.HandleFunc("/registro", userHandler.ShowForm)          // Muestra el HTML
	http.HandleFunc("/procesar-registro", userHandler.SubmitForm) // Guarda en JSONL
	http.HandleFunc("/login", userHandler.ShowForm)
	http.HandleFunc("/procesar-login", userHandler.Login)       // NUEVA: Procesa el inicio de sesión

	// --- ARCHIVOS ESTÁTICOS --- //
	// Importante: Asegúrate de que esta ruta sea exacta
	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Servidor escuchando en http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//go run ./cmd/web/main.go