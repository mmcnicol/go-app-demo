package main

import (
	"log"
	"net/http"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/views"
)

func main() {
	// CRITICAL: The server MUST know the route exists
	//app.Route("/", &Home{})
	//app.Route("/", func() app.Composer { return &Home{} })
	// Register the components that correspond to routes
	//app.Route("/", func() app.Composer { return &pages.LoginPage{} })
	//app.Route("/dev/storybook", func() app.Composer { return &pages.ComponentGallery{} })
	
	//app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	//app.Route("/patients", func() app.Composer { return &pages.PatientList{} })
	//app.Route("/discharge", func() app.Composer { return &pages.DischargeForm{} })
	//app.Route("/discharge/{id}", func() app.Composer { return &pages.DischargeForm{} })

	app.Route("/", func() app.Composer { return &views.BasePage{} })
	app.Route("/dev/components", func() app.Composer { return &view.ComponentGallery{} })
	
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
	
	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
