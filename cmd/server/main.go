package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	// The Handler is the core of go-app server-side.
	// It automatically serves the Wasm binary and static resources.
	appHandler := &app.Handler{
		Name:        "Go-App demo",
		Description: "A minimalist prototype using Go and WebAssembly",
		Author:      "mmcnicol",
		Styles: []string{
			"/web/css/main.css", // Link to your custom CSS
		},
		Icon: app.Icon{
			Default: "/web/images/logo.png",
		},
	}

	// Setup standard HTTP routes
	http.Handle("/", appHandler)
	
	// Example API endpoint
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
