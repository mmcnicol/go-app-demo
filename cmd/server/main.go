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

// Mock API handler for autocomplete - updated to match both formats
func MockAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameter
	query := r.URL.Query().Get("q")
	limit := r.URL.Query().Get("limit")
	
	// Set headers for JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Filter mock data based on query
	var results []map[string]interface{}
	queryLower := strings.ToLower(query)
	
	mockData := []struct {
		id    string
		label string
		value string
	}{
		{"1", "aaaa", "value_aaaa"},
		{"2", "aaaaa", "value_aaaaa"},
		{"3", "aaaaaa", "value_aaaaaa"},
		{"4", "bbbb", "value_bbbb"},
		{"5", "bbbbb", "value_bbbbb"},
		{"6", "bbbbbb", "value_bbbbbb"},
		{"7", "Apple", "fruit_apple"},
		{"8", "Banana", "fruit_banana"},
		{"9", "Cherry", "fruit_cherry"},
		{"10", "Date", "fruit_date"},
	}
	
	// Filter based on query
	for _, item := range mockData {
		if queryLower == "" || strings.Contains(strings.ToLower(item.label), queryLower) {
			results = append(results, map[string]interface{}{
				"id":    item.id,
				"label": item.label,
				"value": item.value,
			})
			if len(results) >= 10 { // Default limit
				break
			}
		}
	}
	
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)
	
	// Create JSON response in structured format
	response := map[string]interface{}{
		"status":  "success",
		"data":    results,
		"query":   query,
		"count":   len(results),
		"limit":   limit,
		"has_more": len(results) >= 10,
	}
	
	json.NewEncoder(w).Encode(response)
}
