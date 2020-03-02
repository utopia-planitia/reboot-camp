package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	roboot "github.com/utopia-planitia/roboot/pkg"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	svc := roboot.NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// https://blog.cloudflare.com/exposing-go-on-the-internet/
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        svc,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on port %s", port)

	return s.ListenAndServe()
}
