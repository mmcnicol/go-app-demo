package pages

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/components"
    "go-app-demo/internal/models"
)

type ComponentGallery struct {
    app.Compo
    selectedComponent string
}

// Define mock data
var mockData = []models.User{
    {Name: "John Doe", Email: "john@example.com"},
    {Name: "Jane Smith", Email: "jane@example.com"},
    {Name: "Bob Johnson", Email: "bob@example.com"},
    {Name: "Alice Williams", Email: "alice@example.com"},
    {Name: "Charlie Brown", Email: "charlie@example.com"},
}

func (g *ComponentGallery) Render() app.UI {
    return app.Div().Class("storybook-layout").Body(
        // Sidebar
        app.Nav().Class("sidebar").Body(
            app.H3().Text("My Components"),
            app.Button().Text("UserTable").OnClick(func(ctx app.Context, e app.Event) {
                g.selectedComponent = "UserTable"
                ctx.Update()
            }),
            app.Button().Text("Navbar").OnClick(func(ctx app.Context, e app.Event) {
                g.selectedComponent = "Navbar"
                ctx.Update()
            }),
        ),
        // Preview Area
        app.Main().Class("preview").Body(
            app.If(g.selectedComponent == "UserTable",
                func() app.UI {
                    //return components.NewUserTable(mockData)
                    return &components.UserTable{
                        Users:       mockData,
                        CurrentPage: 1,
                        PageSize:    10,
                        TotalRows:   len(mockData),
                        SortBy:      "name",
                        SortOrder:   "asc",
                    }
                },
            ).ElseIf(g.selectedComponent == "Navbar",
                func() app.UI {
                    return &components.Navbar{}
                },
            ),
        ),
    )
}
