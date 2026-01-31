package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ApplicationBanner struct {
    app.Compo
}

func NewApplicationBanner() *ApplicationBanner {
    return &ApplicationBanner{}
}

func (a *ApplicationBanner) Render() app.UI {
	return app.Nav().Class("app-banner").Body(
		app.P().Text("Application Banner"),
	)
}
