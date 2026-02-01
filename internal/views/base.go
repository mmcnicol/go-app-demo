package views

import (
	"fmt"
	"time"
	
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
	"go-app-demo/internal/models"
	"go-app-demo/internal/state"
)

type BasePage struct {
	app.Compo
	//navItemsGopherContext    []models.NavItem
	//navItemsNonGopherContext []models.NavItem
	ActiveID                 string // Still used for content routing
	User                     *models.User
	UserSettings             *models.UserSettings
	GopherDemographics       *models.GopherDemographics
	RecentGophers            []models.RecentGopherItem
	LabResults               []models.LabResultItem
	
	// Store navigation component instance
	LeftNavigation    *components.LeftNavigation
	isMounted     bool // Managed internally
}

// OnMount is called when the component is first loaded
func (b *BasePage) OnMount(ctx app.Context) {
	app.Log("[BasePage OnMount] Initializing")

	user := &models.User{
		Username: "iamcheating",
		Forename: "I'm",
		Surname:  "Cheating",
	}

	// a hack: should be based of user's roles & permissions
	userSettings := &models.UserSettings{
		DefaultNavSectionNonGopherContext: "GopherLists",
		DefaultNavSectionGopherContext: "GopherRecords",
		DefaultNavItemNonGopherContext:  "RecentGophers",
		DefaultNavItemGopherContext:  "LabResults",
	}
	
	ctx.SetState(state.UserKey, user)
	ctx.SetState(state.UserSettingsKey, userSettings)

	b.User = user

	ctx.ObserveState(state.UserKey, &b.User).
		OnChange(func() {
			app.Log("[BasePage] User changed")
			b.Refresh(ctx)
		})
	ctx.ObserveState(state.GopherDemographicsKey, &b.GopherDemographics).
		OnChange(func() {
			app.Log("[BasePage] Gopher Demographics changed")
			b.Refresh(ctx)
		})
	ctx.ObserveState(state.ActiveItemKey, &b.ActiveID).
		OnChange(func() {
			app.Log("[BasePage] Active Item changed")
			b.Refresh(ctx)
		})

	// Create gopher navigation instance
	if b.LeftNavigation == nil {
		app.Log("[BasePage Render] Creating new left navigation")
		b.LeftNavigation = components.NewLeftNavigation()
	}
	
	b.isMounted = true
	b.Refresh(ctx)
}

// Refresh updates the component
func (b *BasePage) Refresh(ctx app.Context) {
	ctx.Update()
}

func (b *BasePage) Render() app.UI {
	app.Log("[BasePage Render] Called")
	app.Log("[BasePage Render] Called, isMounted:", b.isMounted)
	
	// Don't render anything until mounted
	if !b.isMounted {
		app.Log("[BasePage Render] Not mounted yet, returning loading/empty")
		//return app.Div().Class("nav-sidebar loading")
		return app.Div()
	}

    var content app.UI
	
	// Check if user is nil (using pointer)
	if b.User == nil {
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
		if b.GopherDemographics == nil {
			// Contextual Routing Logic
			switch b.ActiveID {
			case "RecentGophers":
				recentGophers := b.fetchRecentGophers()
				content = components.NewRecentGophers(recentGophers)
			default:
				content = &components.NotFoundComponent{}
			}
			//b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
			//b.activeID = "RecentGophers"
			
			//app.Log("[BasePage Render] Non-gopher context, gopherDemographics is nil")
			//app.Log("[BasePage Render] Navigation items:", len(b.navItemsNonGopherContext))

			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.LeftNavigation,
				GopherBanner:      nil,
				Body:              content,
				PageFooter:        components.NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
			}
			
		} else {
			// Contextual Routing Logic
			switch b.ActiveID {
			case "RecentGophers":
				b.RecentGophers = b.fetchRecentGophers()
				content = components.NewRecentGophers(b.RecentGophers)
			case "LabResults":
				b.LabResults = b.fetchLabResults()
				content = components.NewLabResults(b.LabResults)
			default:
				content = &components.NotFoundComponent{}
			}
			//b.navItemsGopherContext = b.fetchNavItemsGopherContext()
			//b.activeID = "LabResults"

			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.LeftNavigation,
				GopherBanner:      components.NewGopherBanner(),
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
		
		// a hack: should be based of user's roles & permissions
		userSettings := &models.UserSettings{
			DefaultNavSectionNonGopherContext: "GopherLists",
			DefaultNavSectionGopherContext: "GopherRecords",
			DefaultNavItemNonGopherContext:  "RecentGophers",
			DefaultNavItemGopherContext:  "LabResults",
		}

		ctx.SetState(state.UserKey, user)
		ctx.SetState(state.UserSettingsKey, userSettings)

		ctx.SetState(state.GopherDemographicsKey, nil)
		
		ctx.SetState(state.NavItemsKey, b.fetchNavItemsNonGopherContext())
		
		ctx.SetState(state.ExpandedSectionKey, userSettings.defaultNavSectionNonGopherContext)
		ctx.SetState(state.ActiveItemKey, userSettings.defaultNavItemNonGopherContext)
		
		// Trigger a re-render
		ctx.Update()
		
		return user, nil
	}
	
	return nil, fmt.Errorf("invalid credentials")
}

func (b *BasePage) handleQuickSearch(ctx app.Context, patientID string) error {
	app.Log("[BasePage handleQuickSearch] Called with patientID:", patientID)
	
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
		//app.Log("[BasePage handleQuickSearch async] Starting search simulation")
		
		time.Sleep(1 * time.Second)
		
		newGopherDemographics := b.fetchGophersDemographics()
		//app.Log("[BasePage handleQuickSearch async] Fetched demographics:", newGopherDemographics)
		
		// Update the navigation context
		ctx.Dispatch(func(ctx app.Context) {
			//app.Log("[BasePage handleQuickSearch dispatch] Updating state")
			
			ctx.SetState(state.GopherDemographicsKey, newGopherDemographics)
		
			ctx.SetState(state.NavItemsKey, b.fetchNavItemsGopherContext())
		
			ctx.ObserveState(state.UserSettingsObservable).Value(b.UserSettings)

			ctx.SetState(state.ExpandedSectionObservable, b.UserSettings.defaultNavSectionGopherContext)
			ctx.SetState(state.ActiveItemObservable, b.UserSettings.defaultNavItemNonGopherContext)
			
			//app.Log("[BasePage handleQuickSearch dispatch] Navigation reset for gopher context")
		})
		
		//app.Log("[BasePage] Found patient:", patientID)
	})
	
	return nil
}

// Logout handler in BasePage
func (b *BasePage) handleLogout(ctx app.Context) {
	app.Log("[BasePage] Logging out user")
	
	ctx.SetState(state.UserObservable, nil)
	ctx.SetState(state.UserSettingsObservable, nil)

	ctx.SetState(state.GopherDemographicsObservable, nil)
	
	ctx.SetState(state.NavItemsObservable, nil)
	
	ctx.SetState(state.ExpandedSectionObservable, nil)
	ctx.SetState(state.ActiveItemObservable, nil)

	// Trigger re-render
	ctx.Update()
}

// OnNav is called on Browser Refresh button click or Back button click
func (b *BasePage) OnNav(ctx app.Context) {
	app.Log("[BasePage OnNav] Navigation event")
	// We could check for existing session/cookie here
	// For demo purposes, we'll start with no user logged in

	ctx.SetState(state.UserObservable, nil)
	ctx.SetState(state.UserSettingsObservable, nil)

	ctx.SetState(state.GopherDemographicsObservable, nil)
	
	ctx.SetState(state.NavItemsObservable, nil)
	
	ctx.SetState(state.ExpandedSectionObservable, nil)
	ctx.SetState(state.ActiveItemObservable, nil)

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
			Route:     "/settings",
		},
		{
			ID:                "GopherLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "RecentGophers",
					Label:     "Recent Gophers",
					Icon:      "fas fa-list",
					Route:     "/recent-gophers",
				},
				{
					ID:        "PharmacyDischargeList",
					Label:     "Pharmacy Discharge List",
					Icon:      "fas fa-plus-circle",
					Route:     "/pharmacy-discharge",
				},
			},
		},
		{
			ID:        "GopherSearch",
			Label:     "Gopher Search",
			Icon:      "fas fa-search",
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
			Route:     "/settings",
		},
		{
			ID:                "GopherLists",
			Label:             "Gopher Lists",
			Icon:              "fas fa-user-injured",
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "RecentGophers",
					Label:     "Recent Gophers",
					Icon:      "fas fa-list",
					Route:     "/recent-gophers",
				},
				{
					ID:        "PharmacyDischargeList",
					Label:     "Pharmacy Discharge List",
					Icon:      "fas fa-plus-circle",
					Route:     "/pharmacy-discharge",
				},
			},
		},
		{
			ID:        "GopherSearch",
			Label:     "Gopher Search",
			Icon:      "fas fa-search",
			Route:     "/search",
		},
		{
			ID:                "GopherRecords",
			Label:             "Gopher Records",
			Icon:              "fas fa-user-injured",
			Route:             "",
			Children: []models.NavItem{
				{
					ID:        "LabResults",
					Label:     "Lab Results",
					Icon:      "fas fa-list",
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
