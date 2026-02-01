package state

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/models"
)

var (
	UserObservable = app.Observable{Name: "user"}
	UserSettingsObservable = app.Observable{Name: "userSettings"}
	GopherDemographicsObservable = app.Observable{Name: "gopherDemographics"}
	NavItemsObservable = app.Observable{Name: "navItems"}
    ExpandedSectionObservable = app.Observable{Name: "expandedSection"}
	ActiveItemObservable = app.Observable{Name: "activeNavItem"}
)
