package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// Sidebar component
type sidebar struct {
	app.Compo
}

func (s *sidebar) Render() app.UI {
	return app.Div().Style("width", "200px").Style("float", "left").Body(
		app.Ul().Body(
			app.Li().Body(app.A().Href("/").Text("Dashboard")),
			app.Li().Body(app.A().Href("/users").Text("Users")),
		),
	)
}
