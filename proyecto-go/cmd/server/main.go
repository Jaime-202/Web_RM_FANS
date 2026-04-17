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

	// Servir archivos estáticos y el index.html en la raíz
	// IMPORTANTE: Asegúrate de que la ruta apunta a la carpeta donde está tu index.html original.
	// Por tu captura, veo que index.html y las carpetas css/, img/ y pages/ están fuera de proyecto-go.
	// Como ejecutas desde dentro de proyecto-go, subimos un nivel con "../"
	fs := http.FileServer(http.Dir("../")) 
	mux.Handle("/", fs)

	// Las rutas del formulario se mantienen igual
	mux.HandleFunc("/contacto", contactHandler.ServeForm)     // GET
	mux.HandleFunc("/contacto/enviar", contactHandler.ProcessForm) // POST

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
