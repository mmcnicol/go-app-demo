package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// Usage:
// &PageFooter{}
type PageFooter struct {
	app.Compo
}

func (p *PageFooter) Render() app.UI {
	return app.Div().Class("content-footer").Body(
		app.P().Text("Page Footer"),
	)
}
