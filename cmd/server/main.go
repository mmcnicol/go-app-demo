package main

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/views"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// CRITICAL: The server MUST know the route exists
	// Register the components that correspond to routes
	//app.Route("/", &Home{})
	//app.Route("/", func() app.Composer { return &Home{} })
	//app.Route("/discharge/{id}", func() app.Composer { return &pages.DischargeForm{} })

	app.Route("/", func() app.Composer { return &views.BasePage{} })
	app.Route("/dev/components", func() app.Composer { return &views.ComponentGallery{} })
	
	h := &app.Handler{
		Name:      "Clinical Portal",
		Description: "A minimalist prototype using Go and WebAssembly",
		Author:      "mmcnicol",
		Styles: []string{
			"/web/css/main.css", // Link to your custom CSS
		},
		Icon: app.Icon{
			Default: "/web/images/logo.png",
		},
		//Resources: app.LocalDir("web"),
	}

	http.Handle("/", h)

	// Example API endpoint
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Autocomplete API endpoint
	http.HandleFunc("/api/autocomplete", MockAutocompleteHandler)

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Mock API handler for autocomplete
func MockAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameter
	query := r.URL.Query().Get("q")
	limit := r.URL.Query().Get("limit")
	
	// Set headers for JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Filter mock data based on query
	var results []string
	queryLower := strings.ToLower(query)
	
	mockData := []string{
		"aaaa",
		"aaaaa", 
		"aaaaaa",
		"bbbb",
		"bbbbb",
		"bbbbbb",
		"Apple",
		"Banana",
		"Cherry",
		"Date",
		"Elderberry",
	}
	
	for _, item := range mockData {
		if strings.Contains(strings.ToLower(item), queryLower) {
			results = append(results, item)
			if len(results) >= 10 { // Default limit
				break
			}
		}
	}
	
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)
	
	// Create JSON response
	response := map[string]interface{}{
		"status": "success",
		"data":   results,
		"query":  query,
		"count":  len(results),
		"limit":  limit,
	}
	
	json.NewEncoder(w).Encode(response)
}
