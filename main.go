package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Alphonnse/file_server/handlers"
	"github.com/Alphonnse/file_server/internal/database"
	"github.com/Alphonnse/file_server/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Its a db driver. _ there means include this even if i dnot use is directly
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("DB url is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("can't connect to DB:", err)
	}

	db := database.New(conn)
	apiCfg := handlers.ApiConfig{
		DB: db,
	}

	apiWrapper := middleware.ApiConfigWrapper{
		ApiConfig: apiCfg,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/healthz", handlers.HandlerReadiness)
	router.Get("/err", handlers.ErrorHandler)

	router.Post("/signup", apiCfg.HandlerSignup)
	router.Post("/login", apiCfg.HandlerLogin)

	router.Route("/{username}/{base_dir}", func(r chi.Router) {
		// and than i need to read the dir on OS and than with regexp
		// give him access

		// There might be a trouble with cookie, when using links into site
		r.Get("/{path_of_dirs}*", apiWrapper.MiddlewareAuth(handlers.FS))

		// r.Get("/{path_to_file}/{action}", apiWrapper.MiddlewareAuth(handlers.HandlerView))
		

	})
	// router.Get("/{username}/{base_dir}/{path_to_file}/{action}", apiWrapper.MiddlewareAuth(handlers.HandlerView))
	// router.Get("/{username}/{base_dir}/{path_to_dirs}*", apiWrapper.MiddlewareAuth(handlers.ListFilesNew))
	

	// I need to use url parameters here and username in the start of url
	// router.Get("/disk/upload", handlers.UploadGetHandler)
	// router.Post("/disk/upload", handlers.UploadPostHandler)
	//

	// router.Get("/disk/{some_path}*", handlers.ListFilesOld)

	// router.Get("/disk/view/*", handlers.HandlerView)
	// router.Get("/disk/download/*", handlers.HandlerDownload)

	// тут есть два способа боросться с тем, что идет get от redirecta
	// можно либо созать get обработчик для логина
	// либо сделать, что этот обработывал и get и post
	// ИЛИ ПРОСТО НЕ ЮЗАТЬ ЕГО
	router.Get("/validation", apiWrapper.MiddlewareAuth(handlers.Validate))
	// router.Get("/validation", handlers.Validate) // these how to use queries

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting in port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
