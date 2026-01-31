package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/pages"
	//"go-app-demo/internal/components"
)

/*
type Home struct {
	app.Compo
}

func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Clinical Portal"),
		app.P().Text("Discharge Form ready."),
	)
}
*/

func main() {
	
	// Register the components that correspond to routes
	//app.Route("/", &Home{})
	//app.Route("/", func() app.Composer { return &Home{} })
	app.Route("/", func() app.Composer { return &pages.LoginPage{} })
	app.Route("/dev/storybook", func() app.Composer { return &pages.ComponentGallery{} })
	
	app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	//app.Route("/patients", func() app.Composer { return &pages.PatientList{} })
	//app.Route("/discharge", func() app.Composer { return &pages.DischargeForm{} })
	//app.Route("/discharge/{id}", func() app.Composer { return &pages.DischargeForm{} })
	
	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}
