package main

import (
	"log"
	"net/http"
	"os"


	"github.com/Alphonnse/file_server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)


func main() {
	godotenv.Load(".env")
	
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	router := chi.NewRouter()
	
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:		[]string{"https://*", "http://*"},
		AllowedMethods: 	[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: 	[]string{"*"},
		ExposedHeaders: 	[]string{"Link"},
		AllowCredentials: 	false,
		MaxAge: 			300,
	}))

	router.Get("/healthz", handlers.HandlerReadiness)
	router.Get("/err", handlers.ErrorHandler)

	srv := &http.Server {
		Handler: 	router,
		Addr:		":" + portString,
	}

	log.Printf("Server starting in port %v", portString)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
