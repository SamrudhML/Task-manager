package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	handler := handlers.NewTaskHandler(repo)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/tasks", handler.CreateTask)
	r.Get("/tasks", handler.ListTasks)
	r.Get("/tasks/{id}", handler.GetTask)
	r.Put("/tasks/{id}", handler.UpdateTask)
	r.Delete("/tasks/{id}", handler.DeleteTask)

	addr := envDefault("HTTP_ADDR", ":8080")
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
