package main

import (
	"html/template"
	"log"
	"net/http"

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/handlers"
)

func main() {
	// 1. Configuración de Base de Datos y Handlers
	store := database.NewUserStore("api/users.jsonl")

	tmplForm, err := template.ParseFiles("ui/templates/user_form.html")
	if err != nil {
		log.Fatalf("error cargando template del formulario: %v", err)
	}
	userHandler := handlers.NewUserHandler(tmplForm, store)

	// --- RUTAS --- //

	// RUTA A: Inicio (Atrapa todo lo que sea exactamente "/")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Este es el seguro que comentaste. Lo dejamos activo por seguridad.
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

	// RUTA B: Chaquetas
	http.HandleFunc("/chaquetas", func(w http.ResponseWriter, r *http.Request) {
		tmplChaquetas, err := template.ParseFiles("ui/templates/chaquetas_nina.html")
		if err != nil {
			// Si el archivo HTML no existe o tiene otro nombre (como chaquetas_niña con ñ), saltará este error en la pantalla
			http.Error(w, "Error cargando el archivo HTML de chaquetas", http.StatusInternalServerError)
			return
		}
		tmplChaquetas.Execute(w, nil)
	})

	// RUTA C y D: Registro
	http.HandleFunc("/registro", userHandler.ShowForm)
	http.HandleFunc("/procesar-registro", userHandler.SubmitForm)

	// --- ARCHIVOS ESTÁTICOS --- //
	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- ARRANCAR SERVIDOR --- //
	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
