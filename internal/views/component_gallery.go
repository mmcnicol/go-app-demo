package views

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
	"go-app-demo/internal/models"
)

type ComponentGallery struct {
	app.Compo
	selectedComponent string
}

func (g *ComponentGallery) Render() app.UI {
	return app.Div().Class("storybook-layout").Body(
		// Sidebar
		app.Nav().Class("sidebar").Body(
			app.H3().Text("My Components"),
			
			app.Button().Text("Application Banner").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "ApplicationBanner"
				ctx.Update()
			}),
			app.Button().Text("Page Footer").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "PageFooter"
				ctx.Update()
			}),
			app.Button().Text("Login Form").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "LoginForm"
				ctx.Update()
			}),
			app.Button().Text("Recent Gophers").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "RecentGophers"
				ctx.Update()
			}),
			app.Button().Text("Lab Results").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "LabResults"
				ctx.Update()
			}),
			app.Button().Text("Lab Results (Loading)").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "LabResults-Loading"
				ctx.Update()
			}),
			app.Button().Text("Lab Results (Nil)").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "LabResults-Nil"
				ctx.Update()
			}),
			app.Button().Text("Lab Results (Empty)").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "LabResults-Empty"
				ctx.Update()
			}),
			app.Button().Text("UserTable").OnClick(func(ctx app.Context, e app.Event) {
				g.selectedComponent = "UserTable"
				ctx.Update()
			}),
		),
		// Preview Area
		app.Main().Class("preview").Body(
			g.renderSelectedComponent(),
		),
	)
}





// Then add this helper method:
func (g *ComponentGallery) renderSelectedComponent() app.UI {
    switch g.selectedComponent {        
    case "ApplicationBanner":
        return components.NewApplicationBanner(
            "Component Gallery",
            nil, //g.mockQuickSearchHandler,
            nil //g.mockLogoutHandler,
        )
    
	/*
    case "QuickSearch":
        return components.NewQuickSearch(g.mockQuickSearchHandler)
    */

    case "QuickSearch-Mock":
        return components.NewQuickSearch(func(ctx app.Context, searchTerm string) error {
            app.Log("Mock search for:", searchTerm)
            return nil
        })
        
    case "QuickSearch-Empty":
        return components.NewQuickSearch(nil)
        
    case "PageFooter":
        return components.NewPageFooter("© 2024 Component Gallery")
        
    case "LoginForm":
        mockHandler := func(ctx app.Context, username, password string) (*models.User, error) {
            if username == "demo" && password == "demo123" {
                return &models.User{
                    Username: username,
                    Forename: "Gallery",
                    Surname:  "User",
                }, nil
            }
            return nil, fmt.Errorf("Try: demo/demo123")
        }
        return components.NewLoginForm(mockHandler)
        
    case "RecentGophers":
        return components.NewRecentGophers(mockRecentGophers)
        
    case "LabResults":
        return components.NewLabResults(mockLabResults)
        
    case "LabResults-Loading":
        return &components.LabResults{
            LabResults: []models.LabResultItem{},
            Loading:    true,
        }
        
    case "LabResults-Nil":
        return components.NewLabResults(nil)
        
    case "LabResults-Empty":
        return &components.LabResults{
            LabResults: []models.LabResultItem{},
            Loading:    false,
        }
        
    case "UserTable":
        return &components.UserTable{
            Users:       mockUsers,
            CurrentPage: 1,
            PageSize:    10,
            TotalRows:   len(mockUsers),
            SortBy:      "name",
            SortOrder:   "asc",
        }

    default:
        return app.Div().Class("component-instructions").Body(
            app.H2().Text("Component Gallery"),
            app.P().Text("Select a component from the sidebar to preview it here."),
            app.Ul().Body(
                app.Li().Text("UserTable - Shows a table of users with pagination"),
                app.Li().Text("Navbar - Basic navigation bar"),
                app.Li().Text("ApplicationBanner - Header with quick search"),
                app.Li().Text("QuickSearch - Patient ID search component"),
                app.Li().Text("PageFooter - Footer component"),
                app.Li().Text("LoginForm - Login form with validation"),
                app.Li().Text("RecentGophers - Table of gopher data"),
                app.Li().Text("LabResults - Medical lab results table"),
            ),
        )
    }
}

// Define mock Users
var mockUsers = []models.User{
	{Username: "jdoe", Forename: "John", Surname: "Doe"},
	{Username: "jsmith", Forename: "Jane", Surname: "Smith"},
	{Username: "bjohnson", Forename: "Bob", Surname: "Johnson"},
	{Username: "awilliams", Forename: "Alice", Surname: "Williams"},
	{Username: "cbrown", Forename: "Charlie", Surname: "Brown"},
}

// Define mock RecentGophers
var mockRecentGophers = []models.RecentGopherItem{
	{
		GopherId:         "gopher-001",
		Name:             "Alice",
		DateOfBirth:      time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		DateLastAccessed: time.Now().Add(-24 * time.Hour),
	},
	{
		GopherId:         "gopher-003",
		Name:             "Charlie",
		DateOfBirth:      time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC),
		DateLastAccessed: time.Now().Add(-2 * time.Hour),
	},
	{
		GopherId:         "gopher-002",
		Name:             "Bob",
		DateOfBirth:      time.Date(1995, 10, 5, 0, 0, 0, 0, time.UTC),
		DateLastAccessed: time.Now().Add(-5 * time.Hour),
	},
}

// Define mock LabResults
var mockLabResults = []models.LabResultItem{
	{
		ID:         "rpt-001",
		Subject:    "Blood",
		ReportDate: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:         "rpt-003",
		Subject:    "Pancreas",
		ReportDate: time.Now().Add(-2 * time.Hour),
	},
	{
		ID:         "rpt-002",
		Subject:    "Liver",
		ReportDate: time.Now().Add(-5 * time.Hour),
	},
}
