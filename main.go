package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"solar-system/genesis/db"
	"solar-system/genesis/handler"
	"solar-system/genesis/routes"
	"solar-system/genesis/util"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	loadEnv()

	db := connectDB()
	defer db.Close()

	j := &util.JSON{}

	h := &handler.Handler{DB: db, JSON: j}

	r := chi.NewRouter()

	addCors(r)
	addMiddleware(r)

	routes.MountRoutes(r, h)

	gracefullyServe(r)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		origin := os.Getenv("ORIGIN")
		if len(origin) == 0 {
			log.Fatal("$ORIGIN must be set")
		}
		log.Println(".env file loaded successfully")
	} else {
		origin := os.Getenv("ORIGIN")
		if len(origin) == 0 {
			log.Fatal("$ORIGIN must be set")
		}
		log.Println(".env file loaded successfully")
	}
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("$DSN must be set")
	}
	db, err := db.CreateDB(dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

func addCors(r *chi.Mux) {
	origin := os.Getenv("ORIGIN")
	if origin == "" {
		log.Fatal("$ORIGIN must be set")
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{origin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func addMiddleware(r *chi.Mux) {
	r.Use(middleware.Logger)
}

func gracefullyServe(r *chi.Mux) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port " + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-shutdown
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped cleanly")
}
