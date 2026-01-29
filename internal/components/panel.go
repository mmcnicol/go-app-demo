package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Usage:
// &Panel{ Content: app.Text("Hello World") }
type Panel struct {
	app.Compo
	Content app.UI // This acts like a "slot" or "children"
}

func (w *Panel) Render() app.UI {
	return app.Div().Class("fancy-border").Body(
		app.H3().Text("Wrapper Title"),
		w.Content, // Render whatever was passed in
	)
}
