package pages

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type componentGallery struct {
    app.Compo
    selectedComponent string
}

func (g *componentGallery) Render() app.UI {
    return app.Div().Class("storybook-layout").Body(
        // Sidebar
        app.Nav().Class("sidebar").Body(
            app.H3().Text("My Components"),
            app.Button().Text("UserTable").OnClick(func(ctx app.Context, e app.Event) {
                g.selectedComponent = "UserTable"
            }),
            app.Button().Text("Navbar").OnClick(func(ctx app.Context, e app.Event) {
                g.selectedComponent = "Navbar"
            }),
        ),
        // Preview Area
        app.Main().Class("preview").Body(
            app.If(g.selectedComponent == "UserTable",
                &userTable{users: mockData},
            ).ElseIf(g.selectedComponent == "Navbar",
                &navbar{},
            ),
        ),
    )
}
