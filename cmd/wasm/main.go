package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/pages"
	//"go-app-demo/internal/components"
)

func main() {
	// Register the components that correspond to routes
	app.Route("/", func() app.Composer { return &pages.LoginPage{} })
	app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	app.Route("/dev/storybook", func() app.Composer { return &pages.ComponentGallery{} })

	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}
