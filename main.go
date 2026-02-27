package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello world")

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable is not set")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*","http://*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handleError)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server is running on port %s", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}