package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ApplicationBanner struct {
    app.Compo
	Title       string
}

func NewApplicationBanner(title string) *ApplicationBanner {
    return &ApplicationBanner{
        Title: title,
    }
}

func (a *ApplicationBanner) Render() app.UI {
    return app.Header().Class("app-banner").Body(
        app.Div().Class("banner-left").Body(
            app.H1().Text(a.Title),
        ),
        app.Div().Class("banner-right").Body(
            app.Span().Text("Welcome, " + a.UserName),
        ),
    )
}
