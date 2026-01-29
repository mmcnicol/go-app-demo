package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"go-app-demo/internal/pages"
)

func main() {
	// Register the components that correspond to routes
	app.Route("/", &pages.Home{})
	app.Route("/users", &pages.UserList{})
	app.Route("/dev/storybook", &componentGallery{})

	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}
