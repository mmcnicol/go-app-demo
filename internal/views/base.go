package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
	"go-app-demo/internal/models"
)

/*

# Key Features of this Layout:
* Flex-Shrink Control: The banner and footer are set to stay a fixed height, while the app-body and content-area use flex: 1 to expand and contract based on the browser window size.
* Independent Scrolling: The .main-content has overflow-y: auto, meaning the sidebar stays fixed in place while you scroll through a long medical form.
* Component Architecture: By passing Navigation and Body as app.UI fields, you can swap the content easily as the user clicks different items in the sidebar.

*/

type BasePage struct {
	app.Compo
	navItemsGopherContext []models.NavItem
	navItemsNonGopherContext []models.NavItem
	activeID string
	expandedSecID string
	user models.User
	gopherDemographics models.GopherDemographics
	recentGophers []models.RecentGopherItem
	labResults []models.LabResultItem
}

/*
func (b *BasePage) Render() app.UI {
    // 1. Create the Sidebar
    nav := &LeftNavigation{
        Items:        b.menuItems,
        ActiveItemID: b.activeID,
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
        CurrentPatient:   b.patient,
    }
}
*/

func (b *BasePage) Render() app.UI {
    var content app.UI

	if b.user == nil {
		return &PageLayout{
			applicationBanner: NewApplicationBanner("Gopher Portal"),
			leftNavigation:    nil,
			gopherBanner:      nil,
			Body:              NewLoginForm(),
			pageFooter:        NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
		}
	} else {
		if b.gopherDemographics == nil {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = b.fetchRecentGophers()
				content = NewRecentGophers(b.recentGophers)
			default:
				content = &NotFoundComponent{}
			}
			b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
			b.activeID = "RecentGophers"
			b.expandedSecID = "GophersLists"
			return &PageLayout{
				applicationBanner: NewApplicationBanner("Gopher Portal"),
				leftNavigation:    NewLeftNavigation(b.navItemsNonGopherContext, b.activeID, b.expandedSecID),
				gopherBanner:      nil,
				Body:              content,
				pageFooter:        NewPageFooter("© 2024 Clinical Portal. All rights reserved."),
			}
		} else {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = b.fetchRecentGophers()
				content = NewRecentGophers(b.recentGophers)
			case "LabResults":
				b.labResults = b.fetchLabResults()
				content = NewLabResults(b.labResults)
			default:
				content = &NotFoundComponent{}
			}
			b.gopherDemographics = b.fetchGophersDemographics()
			b.navItemsNonGopherContext = b.fetchNavItemsGopherContext()
			b.activeID = "LabResults"
			b.expandedSecID = "GophersRecords"
			return &PageLayout{
				applicationBanner: NewApplicationBanner("Gopher Portal"),
				leftNavigation:    NewLeftNavigation(b.navItemsPatientContext, b.activeID, b.expandedSecID),
				gopherBanner:      NewGopherBanner(b.gopherDemographics),
				Body:              content,
				pageFooter:        NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
			}
		}
	}
}

// OnMount is called when the component is first loaded
func (b *BasePage) OnMount(ctx app.Context) {
	/*
    // Optional: Initialize with an empty patient or fetch from session
    ctx.SetState("current-patient", Patient{})
	*/

	/*
	d.dischargeForm = models.DischargeForm{
		ID:   "MRN-8829",
		Name: "John Doe",
	}
	*/
}

// OnNav is called on Browser Refresh button click or Back button click
func (b *BasePage) OnNav(ctx app.Context) {
    /*
	// Retrieve the object from state
    ctx.GetState("current-patient", &d.patient)
    d.Update()
	*/
}

func (b *BasePage) loadPatientData(id string) {
    // Simulation: would normally be a database lookup
    if id == "MRN-8829" {
        d.dischargeForm = models.DischargeForm{
			ID:   "MRN-8829",
			Name: "John Doe",
		}
    }
    d.Update()
}

func (b *BasePage) fetchGophersDemographics() []models.GopherDemographics {
	// mock
	return models.GopherDemographics{
		GopherId:         "gopher-001",
		Name:             "Alice",
		DateOfBirth:      time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}

func (b *BasePage) fetchRecentGophers() []models.RecentGopherItem {
	// mock
	return []models.RecentGopherItem{
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
}

func (b *BasePage) fetchLabResults() []models.LabResultItem {
	// mock
	return []models.LabResultItem{
		{
			ID:         "rpt-001",
			Subject:             "Blood",
			//ReportDate:      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			ReportDate: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         "rpt-003",
			Subject:             "Pancreas",
			//ReportDate:      time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			ReportDate: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:         "rpt-002",
			Subject:             "Liver",
			//ReportDate:      time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
			ReportDate: time.Now().Add(-5 * time.Hour),
		},
	}
}

func (b *BasePage) fetchNavItemsNonGopherContext() []models.NavItem {
	return []NavItem{
		{
			ID:    "MySettings",
			Label: "My Settings",
			Icon:  "fas fa-cog",
			IsLoading: false, // Example of eager loading state
		},
		{
			ID:    "GopherLists",
			Label: "Gopher Lists",
			Icon:  "fas fa-user-injured",
			IsDefaultExpanded: true, // This section starts open
			Children: []NavItem{
				{ID: "RecentGophers", Label: "Recent Gophers", Icon: "fas fa-list", IsLoading: true },
				{ID: "PharmacyDischargeList", Label: "Pharmacy Discharge List", Icon: "fas fa-plus-circle", IsLoading: true },
			},
		},
		{
			ID:    "GopherSearch",
			Label: "Gopher Search",
			Icon:  "fas fa-chart-line",
		},
	}
}

func (b *BasePage) fetchNavItemsGopherContext() []models.NavItem {
	return []NavItem{
		{
			ID:    "MySettings",
			Label: "My Settings",
			Icon:  "fas fa-cog",
			IsLoading: false, // Example of eager loading state
		},
		{
			ID:    "GopherLists",
			Label: "Gopher Lists",
			Icon:  "fas fa-user-injured",
			IsDefaultExpanded: false,
			Children: []NavItem{
				{ID: "RecentGophers", Label: "Recent Gophers", Icon: "fas fa-list", IsLoading: true },
				{ID: "PharmacyDischargeList", Label: "Pharmacy Discharge List", Icon: "fas fa-plus-circle", IsLoading: true },
			},
		},
		{
			ID:    "GopherSearch",
			Label: "Gopher Search",
			Icon:  "fas fa-chart-line",
		},
		{
			ID:    "GopherRecords",
			Label: "Gopher Records",
			Icon:  "fas fa-user-injured",
			IsDefaultExpanded: true,
			Children: []NavItem{
				{ID: "LabResults", Label: "Lab Results", Icon: "fas fa-list", IsLoading: true },
			},
		},
	}
}
