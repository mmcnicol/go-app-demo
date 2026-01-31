package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ApplicationBanner struct {
    app.Compo
	Title       string
    QuickSearch *QuickSearch
	OnLogout    func(ctx app.Context)
}

func NewApplicationBanner(title string, onQuickSearch QuickSearchHandler, onLogout func(ctx app.Context)) *ApplicationBanner {
	banner := &ApplicationBanner{
		Title:      title,
		OnLogout:   onLogout,
	}

    // Create quick search component if handler is provided
	if onQuickSearch != nil {
		banner.QuickSearch = NewQuickSearch(onQuickSearch)
	}
	
	return banner
}

// NewSimpleApplicationBanner for backward compatibility
func NewSimpleApplicationBanner(title string) *ApplicationBanner {
	return &ApplicationBanner{
		Title:    title,
	}
}

func (a *ApplicationBanner) Render() app.UI {
    return app.Header().Class("app-banner").Body(
        app.Div().Class("banner-left").Body(
            app.H1().Text(a.Title),
        ),
        app.Div().Class("banner-center").Body(
			// Could add other elements here if needed
		),
        app.Div().Class("banner-right").Body(
			// Quick search component
			app.If(a.QuickSearch != nil,
				func() app.UI {
					return a.QuickSearch
				},
			),
			// User info and logout
			app.Div().Class("user-info").Body(
				app.If(a.OnLogout != nil,
					func() app.UI {
						return app.Button().
							Class("logout-button").
							OnClick(a.onLogoutClick).
							Text("Logout")
					},
				).Else(
					func() app.UI {
						return app.Span().Class("logout-placeholder").Text("logout")
					},
				),
			),
		),
	)
}

func (a *ApplicationBanner) onLogoutClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if a.OnLogout != nil {
		a.OnLogout(ctx)
	}
}
