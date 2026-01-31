package models

type NavItem struct {
	ID                string
	Label             string
	Icon              string    // e.g., "fas fa-user-md"
	IsLoading         bool      // Controls the eager-loading spinner
	IsDefaultExpanded bool      // Configures initial state
	Children          []NavItem // If not empty, this is an expandable section
}
