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
	navItemsPatientContext []NavItem
	navItemsNonPatientContext []NavItem
	activeID string
	user User
	patient Patient
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
		return &MainLayout{
			Navigation: nil,
			Body:       &Login{},
		}
	} else {
		if b.patient == nil {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = fetchRecentGophers()
				content = &RecentGophers{
					recentGophers: b.recentGophers,
					SortBy:      "DateLastAccessed",
					SortOrder:   "desc",
				}
			default:
				content = &NotFoundComponent{}
			}
			return &PageLayout{
				applicationBanner: &ApplicationBanner{}
				leftNavigation: &LeftNavigation{Items: b.navItemsNonPatientContext},
				gopherBanner: nil,
				Body:       content,
				pageFooter: PageFooter(),
			}
		} else {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = fetchRecentGophers()
				content = &RecentGophers{
					recentGophers: b.recentGophers,
					SortBy:      "DateLastAccessed",
					SortOrder:   "desc",
				}
			case "LabResults":
				b.labResults = fetchLabResults()
				content = &LabResults{labResults: b.labResults}
			default:
				content = &NotFoundComponent{}
			}
			return &PageLayout{
				applicationBanner: &ApplicationBanner{}
				leftNavigation: &LeftNavigation{Items: b.navItemsPatientContext, ActiveItemID: b.activeID},
				gopherBanner: GopherBanner{},
				Body:       content,
				pageFooter: PageFooter(),
			}
		}
	}
}

// OnMount is called when the component is first loaded
func (b *BasePage) OnMount(ctx app.Context) {
    // Optional: Initialize with an empty patient or fetch from session
    ctx.SetState("current-patient", Patient{})

	/*
	d.dischargeForm = models.DischargeForm{
		ID:   "MRN-8829",
		Name: "John Doe",
	}
	*/
}

// OnNav is called on Browser Refresh button click or Back button click
func (b *BasePage) OnNav(ctx app.Context) {
    // Retrieve the object from state
    ctx.GetState("current-patient", &d.patient)
    d.Update()
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
