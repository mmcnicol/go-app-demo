package models

type NavItem struct {
	ID                string
	Label             string
	Icon              string    // e.g., "fas fa-user-md"
	Route             string
	IsLoading         bool      // Controls the eager-loading spinner
	IsDefaultExpanded bool      // Configures initial state
	IsDefaultSelected bool
	Children          []NavItem // If not empty, this is an expandable section
}
