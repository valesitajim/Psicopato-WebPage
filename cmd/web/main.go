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

	// 2. Cargamos las plantillas por separado
	tmplLogin, err := template.ParseFiles("ui/templates/login.html")
	if err != nil {
		log.Fatalf("error cargando template de login: %v", err)
	}

	tmplRegister, err := template.ParseFiles("ui/templates/register.html")
	if err != nil {
		log.Fatalf("error cargando template de registro: %v", err)
	}

	// 3. Inicializamos el handler pasando ambas plantillas
	userHandler := handlers.NewUserHandler(tmplLogin, tmplRegister, store)

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
	http.HandleFunc("/procesar-registro", userHandler.SubmitForm) // Guarda en JSONL
	http.HandleFunc("/procesar-login", userHandler.Login)         // Procesa el inicio de sesión
	http.HandleFunc("/login", userHandler.ShowLogin)
	http.HandleFunc("/registro", userHandler.ShowRegister)


		// Usando la raíz del proyecto directamente 
	fs := http.FileServer(http.Dir("ui/static")) // Sin el ./  ni / al final
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Servidor escuchando en http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//go run ./cmd/web/main.go
