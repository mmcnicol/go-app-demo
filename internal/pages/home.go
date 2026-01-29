package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type HomePage struct {
	app.Compo
}

func (p *HomePage) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Welcome to the Prototype"),
		app.P().Text("This is a minimalist Go + Wasm application."),
		app.A().Href("/users").Text("View User Table"),
		app.A().Href("/login").Text("Login Page"),
	)
}
