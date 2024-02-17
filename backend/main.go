package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"demo/router"
)

func main() {

	// Configure the http server
	s := &http.Server{
		Addr:               ":8080",
		Handler:            router.Route(),
		ReadTimeout:        10 * time.Second,	
		WriteTimeout:       10 * time.Second,	
		IdleTimeout:        120 * time.Second,  
		MaxHeaderBytes:     1 << 20,				// The maximum size of the header is 1024 bytes
		ErrorLog:           log.New(os.Stderr, "http: ", log.LstdFlags),
		ReadHeaderTimeout:  10 * time.Second,		// Use to avoid reading the header twice in a row
	}

	// Start the server
	log.Printf("Listening on %s\n", s.Addr)
	s.ListenAndServe()
}
