package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/pages"
	//"go-app-demo/internal/components"
)

func main() {
	// Register the components that correspond to routes
	app.Route("/", &pages.LoginPage{})
	app.Route("/home", &pages.HomePage{})
	app.Route("/dev/storybook", &pages.ComponentGallery{})

	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}
