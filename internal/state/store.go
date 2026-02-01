package state

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    //"go-app-demo/internal/models"
)

// internal/state/store.go
package state

// State keys (observable identifiers)
const (
    UserKey               = "user"
    UserSettingsKey       = "userSettings"
    GopherDemographicsKey = "gopherDemographics"
    NavItemsKey           = "navItems"
    ExpandedSectionKey    = "expandedSection"
    ActiveItemKey         = "activeNavItem"
)

var (
	UserObservable = app.Observable{ Name: UserKey }
	UserSettingsObservable = app.Observable{ Name: UserSettingsKey }
	GopherDemographicsObservable = app.Observable{ Name: GopherDemographicsKey }
	NavItemsObservable = app.Observable{ Name: NavItemsKey }
    ExpandedSectionObservable = app.Observable{ Name: ExpandedSectionKey }
	ActiveItemObservable = app.Observable{ Name: ActiveItemKey }
)
