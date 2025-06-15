package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WaldemarEnns/go-task-manager/domains/task"
)

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Handler      http.Handler
}

func main() {
	s := NewServer()

	httpServer := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s,
	}

	go StartServer(httpServer)

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Default().Println("Shutting down server")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Default().Printf("Error shutting down server: %v", err)
		}
	}()

	wg.Wait()
}

func NewServer() http.Handler {
	mux := http.NewServeMux()

	var handler http.Handler = mux

	RegisterRoutes(mux)

	return handler
}

func StartServer(server *http.Server) {
	log.Default().Println("Starting server on port", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func RegisterRoutes(mux *http.ServeMux) {
	helloWorld := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("GET /")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	health := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("GET /health")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/", helloWorld)
	mux.Handle("/health", health)
	mux.Handle("/tasks", http.HandlerFunc(task.GetTasks))
	mux.Handle("/tasks/{id}", http.HandlerFunc(task.GetTask))
}
