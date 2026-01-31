package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type NotFoundComponent struct {
    app.Compo
}

func (n *NotFoundComponent) Render() app.UI {
    return app.Div().Class("not-found").Body(
        app.P().Text("The requested component could not be found."),
    )
}
