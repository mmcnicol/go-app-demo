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
	activeID                 string
	expandedSecID            string
	user                     *models.User
	gopherDemographics       *models.GopherDemographics
	recentGophers            []models.RecentGopherItem
	labResults               []models.LabResultItem
}

func (b *BasePage) Render() app.UI {
    var content app.UI

	// Check if user is nil (using pointer)
	if b.user == nil {
		// Create login form with callback to BasePage's login method
		loginForm := components.NewLoginForm(b.handleLogin)
		return &components.PageLayout{
			ApplicationBanner: components.NewApplicationBanner("Gopher Portal"),
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
			b.expandedSecID = "GophersLists"
			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner("Gopher Portal"),
				LeftNavigation:    components.NewLeftNavigation(b.navItemsNonGopherContext, b.activeID, b.expandedSecID),
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
			b.gopherDemographics = b.fetchGophersDemographics()
			b.navItemsGopherContext = b.fetchNavItemsGopherContext()
			b.activeID = "LabResults"
			b.expandedSecID = "GophersRecords"
			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner("Gopher Portal"),
				LeftNavigation:    components.NewLeftNavigation(b.navItemsGopherContext, b.activeID, b.expandedSecID),
				GopherBanner:      components.NewGopherBanner(b.gopherDemographics),
				Body:              content,
				PageFooter:        components.NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
			}
		}
	}
}

// handleLogin is the callback function passed to LoginForm
func (b *BasePage) handleLogin(ctx app.Context, username, password string) (*models.User, error) {
	// This is where you would integrate with your authentication system
	// For now, we'll use a simple mock
	
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}
	
	// Simple mock authentication
	if username == "demo" && password == "demo" {
		user := &models.User{
			Username: username,
			Name:     "Demo User",
			Email:    "demo@example.com",
		}
		
		// Set the user on BasePage
		b.user = user
		
		// Initialize other data
		b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
		b.navItemsGopherContext = b.fetchNavItemsGopherContext()
		b.activeID = "RecentGophers"
		b.expandedSecID = "GophersLists"
		
		// Trigger a re-render
		ctx.Update()
		
		return user, nil
	}
	
	return nil, fmt.Errorf("invalid credentials")
}

// OnMount is called when the component is first loaded
func (b *BasePage) OnMount(ctx app.Context) {
	// We could check for existing session/cookie here
	// For demo purposes, we'll start with no user logged in
	b.user = nil
	ctx.Update()
}

// OnNav is called on Browser Refresh button click or Back button click
func (b *BasePage) OnNav(ctx app.Context) {
	// We could check for existing session/cookie here
	// For demo purposes, we'll start with no user logged in
	b.user = nil
    //b.activeID = "RecentGophers"
    ctx.Update()
}

func (b *BasePage) fetchGophersDemographics() *models.GopherDemographics {
	// mock - return pointer
	return &models.GopherDemographics{
		GopherId:    "gopher-001",
		Name:        "Alice",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
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
	return []models.NavItem{
		{
			ID:        "MySettings",
			Label:     "My Settings",
			Icon:      "fas fa-cog",
			IsLoading: false,
			Route:     "/settings",
		},
		{
			ID:                "GophersLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: true,
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
	return []models.NavItem{
		{
			ID:        "MySettings",
			Label:     "My Settings",
			Icon:      "fas fa-cog",
			IsLoading: false,
			Route:     "/settings",
		},
		{
			ID:                "GophersLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: false,
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
			ID:                "GophersRecords",
			Label:             "Gopher Records",
			Icon:              "fas fa-user-injured",
			IsDefaultExpanded: true,
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
