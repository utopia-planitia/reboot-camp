package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	camp "github.com/utopia-planitia/reboot-camp/pkg"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	svc := camp.NewServer()

	// https://blog.cloudflare.com/exposing-go-on-the-internet/
	s := &http.Server{
		Addr:           ":8080",
		Handler:        svc,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listen on :8080")
	err := s.ListenAndServe()

	return err
}
