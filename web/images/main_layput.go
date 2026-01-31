package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainLayout struct {
	app.Compo
	ShowDemographics bool
	CurrentPatient   Patient
	Navigation       app.UI // Injected LeftNavigation
	Body             app.UI // Injected Content (Form/Table)
}

func (l *MainLayout) Render() app.UI {
	return app.Div().Class("app-container").Body(
		// 1. Top Banner
		app.Header().Class("app-banner").Body(
			app.H2().Text("CLINICAL PORTAL v1.0"),
		),

		// 2. Body Split (Sidebar | Content)
		app.Div().Class("app-body").Body(
			// Left Sidebar
			l.Navigation,

			// Right Side Main Area
			app.Main().Class("main-content").Body(
				
				// Optional Demographics Component
				app.If(l.ShowDemographics,
					app.Div().Class("demographics-bar").Body(
						app.Div().Body(
							app.Small().Text("PATIENT NAME"),
							app.P().Style("font-weight", "bold").Text(l.CurrentPatient.Name),
						),
						app.Div().Body(
							app.Small().Text("MRN"),
							app.P().Style("font-weight", "bold").Text(l.CurrentPatient.ID),
						),
					),
				),

				// Content Placeholder (The Form/Data)
				app.Div().Class("content-area").Body(
					l.Body,
				),

				// Footer
				app.Footer().Class("content-footer").Body(
					app.Text("© 2026 Hospital Systems - Session Active"),
				),
			),
		),
	)
}
