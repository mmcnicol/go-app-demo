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
	activeID                 string // Still used for content routing
	//user                     *models.User
	//gopherDemographics       *models.GopherDemographics
	//recentGophers            []models.RecentGopherItem
	//labResults               []models.LabResultItem
	
	// Store navigation component instance
	//leftNavigation    *components.LeftNavigation
}

func (b *BasePage) Render() app.UI {
    var content app.UI

	// temp hack
	user := &models.User{
		Username: "duser2",
		Forename: "Demo",
		Surname:  "User2",
	}
	
	// a hack: should be based of user's roles & permissions
	userSettings := &models.UserSettings{
		defaultNavSectionNonGopherContext: "GopherLists",
		defaultNavSectionGopherContext: "GopherRecords",
		defaultNavItemNonGopherContext:  "RecentGophers",
		defaultNavItemGopherContext:  "LabResults",
	}
	
	ctx.SetState(state.UserObservable, user)
	ctx.SetState(state.UserSettingsObservable, userSettings)

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
		var demographics *models.GopherDemographics
    	ctx.ObserveState(state.GopherDemographicsObservable).Value(&demographics)
		if demographics == nil {
			// Contextual Routing Logic
			switch b.activeID {
			case "RecentGophers":
				recentGophers := b.fetchRecentGophers()
				content = components.NewRecentGophers(recentGophers)
			default:
				content = &components.NotFoundComponent{}
			}
			b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
			b.activeID = "RecentGophers"
			
			app.Log("[BasePage Render] Non-gopher context, gopherDemographics is nil")
			app.Log("[BasePage Render] Navigation items:", len(b.navItemsNonGopherContext))
			
			// Create or reuse non-gopher navigation
			if b.leftNavigation == nil {
				app.Log("[BasePage Render] Creating new non-gopher navigation")
				b.leftNavigation = components.NewLeftNavigation(b.navItemsNonGopherContext)
			} else {
				app.Log("[BasePage Render] Reusing non-gopher navigation")
				// Update items if needed
				b.leftNavigation.SetItems(b.navItemsNonGopherContext)
			}

			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.leftNavigation,
				GopherBanner:      nil,
				Body:              content,
				PageFooter:        components.NewPageFooter("© 2026 Clinical Portal. All rights reserved."),
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
			
			// Create or reuse gopher navigation
			if b.leftNavigation == nil {
				app.Log("[BasePage Render] Creating new gopher navigation")
				b.leftNavigation = components.NewLeftNavigation(b.navItemsGopherContext)
			} else {
				app.Log("[BasePage Render] Reusing gopher navigation")
				// Update items if needed
				b.leftNavigation.SetItems(b.navItemsGopherContext)
			}

			return &components.PageLayout{
				ApplicationBanner: components.NewApplicationBanner(
					"Gopher Portal",
					b.handleQuickSearch, // Quick search handler
					b.handleLogout,      // Logout handler
				),
				LeftNavigation:    b.leftNavigation,
				GopherBanner:      components.NewGopherBanner(demographics),
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
			defaultNavSectionNonGopherContext: "GopherLists",
			defaultNavSectionGopherContext: "GopherRecords",
			defaultNavItemNonGopherContext:  "RecentGophers",
			defaultNavItemGopherContext:  "LabResults",
		}

		ctx.SetState(state.UserObservable, user)
		ctx.SetState(state.UserSettingsObservable, userSettings)

		ctx.SetState(state.GopherDemographicsObservable, nil)
		
		ctx.SetState(state.NavItemsObservable, b.fetchNavItemsNonGopherContext())
		
		ctx.SetState(state.ExpandedSectionObservable, userSettings.defaultNavSectionNonGopherContext)
		ctx.SetState(state.ActiveItemObservable, userSettings.defaultNavItemNonGopherContext)
		
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
			
			ctx.SetState(state.GopherDemographicsObservable, newGopherDemographics)
		
			ctx.SetState(state.NavItemsObservable, b.fetchNavItemsGopherContext())
		
			var userSettings *models.UserSettings
			ctx.ObserveState(state.UserSettingsObservable).Value(&userSettings)

			ctx.SetState(state.ExpandedSectionObservable, userSettings.defaultNavSectionGopherContext)
			ctx.SetState(state.ActiveItemObservable, userSettings.defaultNavItemNonGopherContext)
			
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
		defaultNavSectionNonGopherContext: "GopherLists",
		defaultNavSectionGopherContext: "GopherRecords",
		defaultNavItemNonGopherContext:  "RecentGophers",
		defaultNavItemGopherContext:  "LabResults",
	}
	
	ctx.SetState(state.UserObservable, user)
	ctx.SetState(state.UserSettingsObservable, userSettings)

	b.user = user

	// Initialize navigation items
	b.navItemsNonGopherContext = b.fetchNavItemsNonGopherContext()
	b.navItemsGopherContext = b.fetchNavItemsGopherContext()
	b.activeID = "RecentGophers"
	// Create gopher navigation instance
	if b.leftNavigation == nil {
		app.Log("[BasePage Render] Creating new gopher navigation")
		b.leftNavigation = components.NewLeftNavigation(b.navItemsNonGopherContext)
	}
	app.Log("[BasePage OnMount] leftNavigation initialized")
	app.Log("[BasePage OnMount] User set:", b.user)
	
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
