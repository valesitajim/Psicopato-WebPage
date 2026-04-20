package main

import (
	"html/template"
	"log"
	"net/http" //La biblioteca net/http incluye un servidor HTTP completo listo para producción.

	"TIENDAPATOS/internal/database"
	"TIENDAPATOS/internal/handlers"
)

func main() {
	tmpl, err := template.ParseFiles("web/templates/user_form.html")
	if err != nil {
		log.Fatalf("error cargando templates: %v", err)
	}

	store := database.NewUserStore("api/users.jsonl")
	userHandler := handlers.NewUserHandler(tmpl, store)

	http.HandleFunc("/", userHandler.ShowForm)
	http.HandleFunc("/submit", userHandler.SubmitForm)

	// Archivos estáticos (si añades CSS/JS en web/static)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
