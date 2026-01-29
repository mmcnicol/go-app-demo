package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Navbar struct {
    app.Compo
}

func (n *navbar) Render() app.UI {
	return app.Nav().Style("background-color", "#333").Style("color", "white").Body(
		app.H1().Text("My App Admin"),
	)
}
