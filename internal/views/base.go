package views

import (
	"fmt"
	"time"
	
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
	"go-app-demo/internal/models"
)

type BasePage struct {
	app.Compo
	navItemsGopherContext    []models.NavItem
	navItemsNonGopherContext []models.NavItem
	activeID                 string // Still used for content routing
	user                     *models.User
	gopherDemographics       *models.GopherDemographics
	recentGophers            []models.RecentGopherItem
	labResults               []models.LabResultItem
	leftNavigation           *components.LeftNavigation // Single instance
}

func (b *BasePage) Render() app.UI {
    var content app.UI

	// temp hack
	user := &models.User{
		Username: "duser2",
		Forename: "Demo",
		Surname:  "User2",
	}
	
	// Set the user on BasePage
	b.user = user
	
	// Check if user is nil (using pointer)
	if b.user == nil {
		// Create login form with callback to BasePage's login method
		loginForm := components.NewLoginForm(b.handleLogin)
		return &components.PageLayout{
			ApplicationBanner: components.NewSimpleApplicationBanner("Gopher Portal"),
			LeftNavigation:    nil,
			GopherBanner:      nil,
			Body:              loginForm,
			PageFooter:        components.NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
		}
	} else {
		if b.gopherDemographics == nil {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = b.fetchRecentGophers()
				content = components.NewRecentGophers(b.recentGophers)
			default:
				content = &components.NotFoundComponent{}
			}
			b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
			b.activeID = "RecentGophers"
			
			app.Log("[BasePage Render] Non-gopher context, gopherDemographics is nil")
			app.Log("[BasePage Render] Navigation items:", len(b.navItemsNonGopherContext))
			
			// Create or update left navigation
			if b.leftNavigation == nil {
				app.Log("[BasePage Render] Creating new LeftNavigation for non-gopher context")
				b.leftNavigation = components.NewLeftNavigation(b.navItemsNonGopherContext)
			} else {
				app.Log("[BasePage Render] Updating existing LeftNavigation for non-gopher context")
				b.leftNavigation.SetItems(b.Context(), b.navItemsNonGopherContext)
			}
			
			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.leftNavigation, // Use the same instance
				GopherBanner:      nil,
				Body:              content,
				PageFooter:        components.NewPageFooter("© 2024 Clinical Portal. All rights reserved."),
			}
		} else {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				b.recentGophers = b.fetchRecentGophers()
				content = components.NewRecentGophers(b.recentGophers)
			case "LabResults":
				b.labResults = b.fetchLabResults()
				content = components.NewLabResults(b.labResults)
			default:
				content = &components.NotFoundComponent{}
			}
			b.navItemsGopherContext = b.fetchNavItemsGopherContext()
			b.activeID = "LabResults"
			
			app.Log("[BasePage Render] Gopher context, gopherDemographics exists")
			app.Log("[BasePage Render] Navigation items (gopher context):", len(b.navItemsGopherContext))
			app.Log("[BasePage Render] Active ID:", b.activeID)
			
			// Create or update left navigation
			if b.leftNavigation == nil {
				app.Log("[BasePage Render] Creating new LeftNavigation for gopher context")
				b.leftNavigation = components.NewLeftNavigation(b.navItemsGopherContext)
			} else {
				app.Log("[BasePage Render] Updating existing LeftNavigation for gopher context")
				b.leftNavigation.SetItems(b.Context(), b.navItemsGopherContext)
				// Optionally set LabResults as active
				b.leftNavigation.SetActiveItem(b.Context(), "LabResults")
			}
			
			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.leftNavigation, // Use the same instance
				GopherBanner:      components.NewGopherBanner(b.gopherDemographics),
				Body:              content,
				PageFooter:        components.NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
			}
		}
	}
}

// handleLogin is the callback function passed to LoginForm
func (b *BasePage) handleLogin(ctx app.Context, username string, password string) (*models.User, error) {
	// This is where you would integrate with your authentication system
	// For now, we'll use a simple mock
	
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}
	
	// Simple mock authentication
	if username == "demo" && password == "demo123" {
		user := &models.User{
			Username: username,
			Forename: "Demo",
			Surname:  "User",
		}
		
		// Set the user on BasePage
		b.user = user
		
		// Initialize other data
		b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
		b.navItemsGopherContext = b.fetchNavItemsGopherContext()
		b.activeID = "RecentGophers"
		
		// Trigger a re-render
		ctx.Update()
		
		return user, nil
	}
	
	return nil, fmt.Errorf("invalid credentials")
}

// Quick search handler in BasePage
func (b *BasePage) handleQuickSearch(ctx app.Context, patientID string) error {
	app.Log("[BasePage handleQuickSearch] Called with patientID:", patientID)
	app.Log("[BasePage handleQuickSearch] Current gopherDemographics:", b.gopherDemographics)
	app.Log("[BasePage handleQuickSearch] Current activeID:", b.activeID)
	
	if len(patientID) != 10 {
		return fmt.Errorf("Patient ID must be exactly 10 digits")
	}
	
	// Validate numeric
	for _, char := range patientID {
		if char < '0' || char > '9' {
			return fmt.Errorf("Patient ID must contain only numbers")
		}
	}
	
	app.Log("[BasePage] Quick searching for patient:", patientID)
	
	// For demo, just show a message
	ctx.Async(func() {
		app.Log("[BasePage handleQuickSearch async] Starting search simulation")
		
		time.Sleep(1 * time.Second)
		
		newGopherDemographics := b.fetchGophersDemographics()
		app.Log("[BasePage handleQuickSearch async] Fetched demographics:", newGopherDemographics)
		
		// Update the navigation context
		ctx.Dispatch(func(ctx app.Context) {
			app.Log("[BasePage handleQuickSearch dispatch] Updating state")
			app.Log("[BasePage handleQuickSearch dispatch] Old gopherDemographics:", b.gopherDemographics)
			app.Log("[BasePage handleQuickSearch dispatch] Old activeID:", b.activeID)
			
			b.gopherDemographics = newGopherDemographics
			b.navItemsGopherContext = b.fetchNavItemsGopherContext()
			b.activeID = "LabResults" // Set default active item for gopher context
			
			app.Log("[BasePage handleQuickSearch dispatch] New gopherDemographics:", b.gopherDemographics)
			app.Log("[BasePage handleQuickSearch dispatch] New activeID:", b.activeID)
			app.Log("[BasePage handleQuickSearch dispatch] Navigation items updated to gopher context")
			
			// Update left navigation if it exists
			if b.leftNavigation != nil {
				b.leftNavigation.SetItems(ctx, b.navItemsGopherContext)
				b.leftNavigation.SetActiveItem(ctx, "LabResults")
			}
		})
		
		app.Log("[BasePage] Found patient:", patientID)
	})
	
	return nil
}

// Logout handler in BasePage
func (b *BasePage) handleLogout(ctx app.Context) {
	app.Log("[BasePage] Logging out user")
	
	// Clear user session
	b.user = nil
	
	// Clear other state
	b.gopherDemographics = nil
	b.recentGophers = nil
	b.labResults = nil
	b.leftNavigation = nil // Reset navigation instance
	
	// Trigger re-render
	ctx.Update()
}

// OnMount is called when the component is first loaded
func (b *BasePage) OnMount(ctx app.Context) {
	app.Log("[BasePage OnMount] Initializing")
	
	b.user = &models.User{
		Username: "iamcheating",
		Forename: "I'm",
		Surname:  "Cheating",
	}

	// Initialize navigation items
	b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
	b.navItemsGopherContext = b.fetchNavItemsGopherContext()
	b.activeID = "RecentGophers"
	
	app.Log("[BasePage OnMount] Navigation initialized")
	app.Log("[BasePage OnMount] User set:", b.user)
	
	ctx.Update()
}

// OnNav is called on Browser Refresh button click or Back button click
func (b *BasePage) OnNav(ctx app.Context) {
	app.Log("[BasePage OnNav] Navigation event")
	// We could check for existing session/cookie here
	// For demo purposes, we'll start with no user logged in
	b.user = nil
	b.leftNavigation = nil // Reset navigation instance
    ctx.Update()
}

func (b *BasePage) fetchGophersDemographics() *models.GopherDemographics {
	app.Log("[BasePage] fetchGophersDemographics called")
	// mock - return pointer
	return &models.GopherDemographics{
		GopherId:    "gopher-001",
		Name:        "Alice",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}

func (b *BasePage) fetchRecentGophers() []models.RecentGopherItem {
	app.Log("[BasePage] fetchRecentGophers called")
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
	app.Log("[BasePage] fetchLabResults called")
	// mock
	return []models.LabResultItem{
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
}

func (b *BasePage) fetchNavItemsNonGopherContext() []models.NavItem {
	app.Log("[BasePage] fetchNavItemsNonGopherContext called")
	return []models.NavItem{
		{
			ID:        "MySettings",
			Label:     "My Settings",
			Icon:      "fas fa-cog",
			IsLoading: false,
			Route:     "/settings",
		},
		{
			ID:                "GopherLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: true, // This section should be expanded by default
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "RecentGophers",
					Label:     "Recent Gophers",
					Icon:      "fas fa-list",
					IsLoading: false,
					Route:     "/recent-gophers",
				},
				{
					ID:        "PharmacyDischargeList",
					Label:     "Pharmacy Discharge List",
					Icon:      "fas fa-plus-circle",
					IsLoading: false,
					Route:     "/pharmacy-discharge",
				},
			},
		},
		{
			ID:        "GopherSearch",
			Label:     "Gopher Search",
			Icon:      "fas fa-search",
			IsLoading: false,
			Route:     "/search",
		},
	}
}

func (b *BasePage) fetchNavItemsGopherContext() []models.NavItem {
	app.Log("[BasePage] fetchNavItemsGopherContext called")
	return []models.NavItem{
		{
			ID:        "MySettings",
			Label:     "My Settings",
			Icon:      "fas fa-cog",
			IsLoading: false,
			Route:     "/settings",
		},
		{
			ID:                "GopherLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: false, // Not expanded by default in gopher context
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "RecentGophers",
					Label:     "Recent Gophers",
					Icon:      "fas fa-list",
					IsLoading: false,
					Route:     "/recent-gophers",
				},
				{
					ID:        "PharmacyDischargeList",
					Label:     "Pharmacy Discharge List",
					Icon:      "fas fa-plus-circle",
					IsLoading: false,
					Route:     "/pharmacy-discharge",
				},
			},
		},
		{
			ID:        "GopherSearch",
			Label:     "Gopher Search",
			Icon:      "fas fa-search",
			IsLoading: false,
			Route:     "/search",
		},
		{
			ID:                "GopherRecords",
			Label:             "Gopher Records",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: true, // This section should be expanded by default in gopher context
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "LabResults",
					Label:     "Lab Results",
					Icon:      "fas fa-list",
					IsLoading: false,
					Route:     "/lab-results",
				},
			},
		},
	}
}

// Login method (example)
func (b *BasePage) Login(username, password string) bool {
    /*
	// Mock login - replace with actual authentication
    if username == "demo" && password == "demo" {
        b.user = &models.User{
            ID:    "user-001",
            Name:  "Demo User",
            Email: "demo@example.com",
        }
        //ctx.Update()
        return true
    }
	*/
    return false
}

// Logout method
func (b *BasePage) Logout() {
    b.user = nil
    //ctx.Update()
}
