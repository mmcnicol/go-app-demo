package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

/*
type HomePage struct {
	app.Compo
}

func (p *HomePage) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Welcome to the Prototype"),
		app.P().Text("This is a minimalist Go + Wasm application."),
		app.A().Href("/users").Text("View User Table"),
		app.A().Href("/login").Text("Login Page"),
	)
}
*/


/*

# Key Features of this Layout:
* Flex-Shrink Control: The banner and footer are set to stay a fixed height, while the app-body and content-area use flex: 1 to expand and contract based on the browser window size.
* Independent Scrolling: The .main-content has overflow-y: auto, meaning the sidebar stays fixed in place while you scroll through a long medical form.
* Component Architecture: By passing Navigation and Body as app.UI fields, you can swap the content easily as the user clicks different items in the sidebar.

*/

type HomePage struct {
	app.Compo
	menuItems []NavItem
	activeID string
	patient Patient
}

func (h *Home) Render() app.UI {
    // 1. Create the Sidebar
    nav := &LeftNavigation{
        Items:        h.menuItems,
        ActiveItemID: h.activeID,
    }

    // 2. Create the Form/Content
    content := app.Div().Body(
        app.H3().Text("Patient Discharge Form"),
        // ... your form inputs here ...
    )

    // 3. Assemble the Page
    return &MainLayout{
        Navigation:       nav,
        Body:             content,
        ShowDemographics: true,
        CurrentPatient:   h.patient,
    }
}
