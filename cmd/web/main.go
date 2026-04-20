package main

import (
	"html/template"
	"log"
	"net/http"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/handlers"
)

func main() {
	// 1. Configurar la base de datos
	store := database.NewUserStore("api/users.jsonl")

	// 2. Cargar el template del formulario para tu handler
	tmplForm, err := template.ParseFiles("ui/templates/user_form.html")
	if err != nil {
		log.Fatalf("error cargando template del formulario: %v", err)
	}
	userHandler := handlers.NewUserHandler(tmplForm, store)


	// --- CONFIGURACIÓN DE RUTAS --- //

	// RUTA A: La página principal de tu tienda (el index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Evitar que otras rutas raras caigan aquí por error
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmplIndex, err := template.ParseFiles("ui/templates/index.html")
		if err != nil {
			http.Error(w, "Error cargando la página de inicio", http.StatusInternalServerError)
			return
		}
		tmplIndex.Execute(w, nil)
	})

	// RUTA B: El formulario de registro
	// Al entrar en http://localhost:8080/registro se mostrará el user_form.html
	http.HandleFunc("/registro", userHandler.ShowForm)

	// RUTA C: Procesar el formulario
	// Cuando el usuario haga clic en "Crear cuenta", los datos irán a esta ruta
	http.HandleFunc("/procesar-registro", userHandler.SubmitForm)


	// --- ARCHIVOS ESTÁTICOS (CSS y AVIF) --- //
	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


	// --- ARRANCAR SERVIDOR --- //
	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}

	// RUTA D: Página de Chaquetas de niña
	http.HandleFunc("/chaquetas", func(w http.ResponseWriter, r *http.Request) {
		tmplChaquetas, err := template.ParseFiles("ui/templates/chaquetas_niña.html")
		if err != nil {
			http.Error(w, "Error cargando la página de chaquetas", http.StatusInternalServerError)
			return
		}
		tmplChaquetas.Execute(w, nil)
	})
}


