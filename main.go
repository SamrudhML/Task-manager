package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"task-manager/internals/auth"
	"time"

	"task-manager/internals/handlers"
	"task-manager/internals/repositories"

	"github.com/go-chi/chi"
)

func envDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	mongoHost := envDefault("MONGO_HOST", "localhost")
	mongoPort := envDefault("MONGO_PORT", "27017")
	dbName := envDefault("MONGO_DB", "taskdb")
	uri := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := repositories.NewMongoClient(ctx, uri)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	repo := repositories.NewTaskRepository(client, dbName, "tasks")
	userRepo := repositories.NewUserRepository(client, dbName, "users")
	handler := handlers.NewTaskHandler(repo)
	jwtSecret := envDefault("JWT_SECRET", "replace-me-in-production")
	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)

	r := chi.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(pr chi.Router) {
		pr.Use(auth.Middleware(jwtSecret))
		pr.Post("/tasks", handler.CreateTask)
		pr.Get("/tasks", handler.ListTasks)
		pr.Get("/tasks/{id}", handler.GetTask)
		pr.Put("/tasks/{id}", handler.UpdateTask)
		pr.Delete("/tasks/{id}", handler.DeleteTask)
	})

	addr := envDefault("HTTP_ADDR", ":8080")
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
