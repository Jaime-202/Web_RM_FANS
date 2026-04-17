package main

import (
	"log"
	"net/http"
	"os"
	"proyecto-go/internal/database"
	"proyecto-go/internal/handlers"
	"proyecto-go/internal/services"
)

func main() {
	// Asegúrate de que existe la carpeta para los datos
	os.MkdirAll("data", os.ModePerm)

	// 1. Inicializar Capa de Datos
	contactRepo := database.NewJSONContactRepository("data/contacts.jsonline")

	// 2. Inicializar Capa de Lógica de Negocio
	contactService := services.NewContactService(contactRepo)

	// 3. Inicializar Capa de Presentación
	contactHandler := handlers.NewContactHandler(contactService)

	// 4. Configurar el multiplexor (Router) de HTTP
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/contacto", contactHandler.ServeForm)     // GET
	mux.HandleFunc("/contacto/enviar", contactHandler.ProcessForm) // POST

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
