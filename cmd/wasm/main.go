package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/pages"
	//"go-app-demo/internal/components"
)

func main() {
	// Register the components that correspond to routes
	app.Route("/", func() app.UI { return &pages.LoginPage{} } )
	app.Route("/home", func() app.UI { return &pages.HomePage{} } )
	app.Route("/dev/storybook", func() app.UI { return &pages.ComponentGallery{} } )

	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}
