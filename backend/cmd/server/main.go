package main

import (
	"log"
	"net/http"
	"os"

	"go-links/internal/db"
	"go-links/internal/handlers"
	"go-links/internal/repository"
	"go-links/internal/service"
)

func main() {
	dbPath := os.Getenv("GOLINKS_DB_PATH")
	if dbPath == "" {
		dbPath = "golinks.db"
	}

	con, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer con.Close()

	repo := repository.NewLinkRepository(con)
	svc := service.NewLinkService(repo)
	handler := handlers.NewLinkHandler(svc)
	router := handlers.NewRouter(handler)

	addr := os.Getenv("GOLINKS_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("go-links server listening on %s (db: %s)", addr, dbPath)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
