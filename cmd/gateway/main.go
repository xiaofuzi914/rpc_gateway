package main

import (
	"log"
	"net/http"

	"github.com/xiaofuzi914/rpc-gateway/internal/config"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Print loaded configuration for debugging
	log.Printf("Loaded config: %+v", cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Health check requested")
		// w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("listening on", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
