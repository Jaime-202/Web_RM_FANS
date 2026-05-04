package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"proyecto-go/internal/database"
	"proyecto-go/internal/handlers"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/services"
)

func main() {
	// Asegúrate de que existe la carpeta para los datos
	os.MkdirAll("data", os.ModePerm)

	// 1. Inicializar Capa de Datos (Repositorios)
	contactRepo := database.NewJSONContactRepository("data/contacts.jsonline")
	userRepo := database.NewJSONUserRepository("data/users.json")
	newsRepo := database.NewJSONNewsRepository("data/news.json")
	ratingRepo := database.NewJSONRatingRepository("data/ratings.json")
	playerRepo := database.NewJSONPlayerRepository("data/players.json")
	transferRepo := database.NewJSONTransferRepository("data/transfers.json")

	// 2. Inicializar Capa de Lógica de Negocio
	contactService := services.NewContactService(contactRepo)

	// 3. Inicializar Plantillas (Templates)
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	// 4. Inicializar Capa de Presentación (Handlers)
	contactHandler := handlers.NewContactHandler(contactService, tmpl)
	authHandler := handlers.NewAuthHandler(userRepo, tmpl)
	newsHandler := handlers.NewNewsHandler(newsRepo, ratingRepo, tmpl)
	playerHandler := handlers.NewPlayerHandler(playerRepo, tmpl)
	transferHandler := handlers.NewTransferHandler(transferRepo, tmpl)
	homeHandler := handlers.NewHomeHandler(newsRepo, playerRepo, transferRepo, tmpl)

	// 5. Configurar el multiplexor (Router) de HTTP
	mux := http.NewServeMux()

	// Rutas Públicas (Sin Auth obligatoria)
	mux.HandleFunc("/login", authHandler.ServeLogin)
	mux.HandleFunc("/register", authHandler.ServeRegister)
	mux.HandleFunc("/logout", authHandler.ServeLogout)

	// Rutas Base (aplicamos AuthMiddleware para saber si está logueado pero sin bloquear el acceso de lectura)
	mux.Handle("/", middleware.AuthMiddleware(http.HandlerFunc(homeHandler.ServeHome)))
	mux.Handle("/noticias", middleware.AuthMiddleware(http.HandlerFunc(newsHandler.ServeList)))
	mux.Handle("/noticia", middleware.AuthMiddleware(http.HandlerFunc(newsHandler.ServeDetail)))
	mux.Handle("/plantilla", middleware.AuthMiddleware(http.HandlerFunc(playerHandler.ServeList)))
	mux.Handle("/fichajes", middleware.AuthMiddleware(http.HandlerFunc(transferHandler.ServeList)))

	// Formulario de contacto
	mux.Handle("/contacto", middleware.AuthMiddleware(http.HandlerFunc(contactHandler.ServeForm)))
	mux.Handle("/contacto/enviar", middleware.AuthMiddleware(http.HandlerFunc(contactHandler.ProcessForm)))

	// Rutas de Usuario Normal (Requieren estar logueado)
	mux.Handle("/news/rate", middleware.AuthMiddleware(middleware.RequireAuth(newsHandler.ProcessRating)))

	// Rutas de Administrador (Requieren rol admin)
	mux.Handle("/admin/news/create", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(newsHandler.ServeCreateForm))))
	mux.Handle("/admin/news/edit", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(newsHandler.ServeEditForm))))
	mux.Handle("/admin/news/delete", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(newsHandler.ServeDelete))))
	
	mux.Handle("/admin/player/create", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(playerHandler.ServeCreateForm))))
	mux.Handle("/admin/player/edit", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(playerHandler.ServeEditForm))))
	mux.Handle("/admin/player/delete", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(playerHandler.ServeDelete))))
	
	// API REST para CRUD de Jugadores
	mux.Handle("/api/players", middleware.AuthMiddleware(http.HandlerFunc(playerHandler.ServeAPI)))

	mux.Handle("/admin/transfer/create", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(transferHandler.ServeCreateForm))))
	mux.Handle("/admin/transfer/edit", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(transferHandler.ServeEditForm))))
	mux.Handle("/admin/transfer/delete", middleware.AuthMiddleware(middleware.RequireAuth(middleware.RequireAdmin(transferHandler.ServeDelete))))

	// API REST para CRUD de Fichajes
	mux.Handle("/api/transfers", middleware.AuthMiddleware(http.HandlerFunc(transferHandler.ServeAPI)))

	// Archivos estáticos (Solo css, img y js)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../css"))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("../img"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../js"))))

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
