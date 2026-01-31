package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
)

/*

# Requirements

LeftNavigation component:
* clicking to expand a section should collapse any other expanded sections.
* "Accordion" logic to ensure only one section stays open at a time.
* eager-loading spinner state.
* state of what is "active" (selected)
* state of what is "expanded" (the accordion).


# Data Structure (The Model)

To handle clickable items versus expandable sections, we’ll use a nested struct. This allows you to define your entire menu in one place.

NavItem Struct:
* ID: Unique identifier (used for selection and expansion state).
* Label: Display text (e.g., "Patient Records").
* Icon: The name of the icon (e.g., "fa-user").
* IsLoading: Boolean state for your eager-loading feature.
* Items: A slice of NavItem (if this slice is not empty, the item behaves as an Expandable Section).
* IsDefaultExpanded: Boolean for initial state.


# Visual Design (The UI/UX)

For a clinical environment, high contrast and clarity are key.

Color Palette:
* Background: Dark Slate or Deep Charcoal (#1e293b).
* Text: Stark White (#ffffff) for active items, Muted Gray (#94a3b8) for inactive.
* Accent: "Medical Blue" (#3b82f6) for the selection indicator.
* Selected State: A vertical blue bar on the left edge or a subtle background highlight.
* Loading State: Replace the static icon with a CSS spin animation when IsLoading is true.


# Behavioral Logic (The "Brain")

The LeftNavigation component will manage its own internal state to handle the "Accordion" behavior.
* Selection Logic: A CurrentSelection string matches the NavItem.ID.
* Accordion Logic: A CurrentExpandedSection string. When a user clicks a section:
  * If the section is already open, it collapses.
  * If it is closed, it opens and automatically clears the previous CurrentExpandedSection.
* Eager Loading: The component renders a spinner icon conditionally based on the IsLoading property of the specific item being clicked.


# Why this works for a Clinical App

* Fitts's Law: The 4px blue border-left provides a massive visual anchor, helping clinicians quickly identify their location in a complex menu.
* Cognitive Load: The transition from muted gray to bright white/blue on hover reduces "noise" so the user can focus on the active task.
* Responsiveness: The CSS spinner is purely client-side, ensuring that even if the WASM logic is busy "eager loading" data, the UI provides immediate feedback that something is happening.


# Why this design works:
* State-Driven Accordion: By using a single ExpandedSecID string, the "collapse others" logic is inherent—as soon as you assign a new ID to that variable, the app.If condition for the previous section becomes false, and it disappears from the DOM.
* Recursive Rendering: The renderItem function calls itself for children, allowing you to have multiple levels of nesting if your clinical portal ever grows in complexity.
* Eager Loading Ready: The IsLoading check is baked into the icon container. When you toggle that boolean in your Go code, the spinner will instantly replace the icon.


# Usage

```Go
menu := []NavItem{
    {
        ID:    "dashboard",
        Label: "Dashboard",
        Icon:  "fas fa-chart-line",
    },
    {
        ID:    "patients",
        Label: "Patient Records",
        Icon:  "fas fa-user-injured",
        IsDefaultExpanded: true, // This section starts open
        Children: []NavItem{
            {ID: "p-list", Label: "Active Patients", Icon: "fas fa-list"},
            {ID: "p-admit", Label: "Admissions", Icon: "fas fa-plus-circle"},
        },
    },
    {
        ID:    "settings",
        Label: "System Settings",
        Icon:  "fas fa-cog",
        IsLoading: true, // Example of eager loading state
    },
}

// In your main layout:
// &LeftNavigation{Items: menu}
```
*/

type LeftNavigation struct {
	app.Compo
	Items         []models.NavItem
	ActiveItemID  string // Tracks the currently selected page
	ExpandedSecID string // Tracks the currently open accordion section
}

func NewLeftNavigation(items []models.NavItem, activeItemID string, expandedSecID string) *LeftNavigation {
    return &LeftNavigation{
        Items:         items,
        ActiveItemID:  activeItemID,
        ExpandedSecID: expandedSecID,
    }
}

// OnMount handles the "Default Expanded" requirement
func (n *LeftNavigation) OnMount(ctx app.Context) {
	for _, item := range n.Items {
		if item.IsDefaultExpanded && len(item.Children) > 0 {
			n.ExpandedSecID = item.ID
			break
		}
	}
}

// Helper method to recursively find and set active item
func (n *LeftNavigation) findAndSetActiveItem(path string, items []models.NavItem) bool {
    for _, item := range items {
        if item.ID == path || item.Route == path { // Check both ID and Route if exists
            n.ActiveItemID = item.ID
            return true
        }
        // Check children recursively
        if len(item.Children) > 0 {
            if n.findAndSetActiveItem(path, item.Children) {
                // Also expand parent if child is active
                n.ExpandedSecID = item.ID
                return true
            }
        }
    }
    return false
}

/* 
To ensure the sidebar knows which item is active when a user 
refreshes the page or uses the back button, use app.Window().URL() 
inside the OnMount or OnNav lifecycle events.
*/
func (n *LeftNavigation) OnNav(ctx app.Context) {
    // Automatically highlight the item that matches the current URL
    currPath := ctx.Page().URL().Path
    
	/*
	n.ActiveItemID = currPath
    ctx.Update()
	*/

	// Find matching item in the navigation tree
    n.findAndSetActiveItem(currPath, n.Items)
    
    ctx.Update()
}

func (n *LeftNavigation) Render() app.UI {
	return app.Nav().Class("nav-sidebar").Body(
		app.Ul().Class("nav-list").Body(
			app.Range(n.Items).Map(func(i int) app.UI {
				return n.renderItem(n.Items[i], false)
			}),
		),
	)
}

// renderItem creates the UI for an individual item or a section
func (n *LeftNavigation) renderItem(item models.NavItem, isChild bool) app.UI {
	hasChildren := len(item.Children) > 0
	isExpanded := n.ExpandedSecID == item.ID
	isSelected := n.ActiveItemID == item.ID

	return app.Li().Body(
		app.Div().
			Class("nav-item").
			// Apply conditional classes for selection and children
			Class(app.If(isSelected, "selected").Else("")).
			Class(app.If(hasChildren, "has-children").Else("")).
			OnClick(func(ctx app.Context, e app.Event) {
				n.handleItemClick(ctx, item)
			}).
			Body(
				// Icon or Spinner (Eager Loading Logic)
				app.Div().Class("nav-icon").Body(
					app.If(item.IsLoading,
						app.Span().Class("spinner"),
					).Else(
						app.I().Class(item.Icon),
					),
				),
				//app.Text(item.Label),
				app.Span().Class("nav-label").Text(item.Label),
				// Chevron indicator for sections
				app.If(hasChildren,
					app.I().Class("fas ml-auto").
						Class(app.If(isExpanded, "fa-chevron-down").Else("fa-chevron-right")),
				),
			),
		// Submenu (The Accordion Content)
		app.If(hasChildren && isExpanded,
			app.Ul().Class("nav-submenu").Body(
				app.Range(item.Children).Map(func(j int) app.UI {
					return n.renderItem(item.Children[j], true)
				}),
			),
		),
	)
}

func (n *LeftNavigation) handleItemClick(ctx app.Context, item models.NavItem) {
	// 1. If it's a section (has children), toggle accordion
	if len(item.Children) > 0 {
		if n.ExpandedSecID == item.ID {
			// If already open, close it
			n.ExpandedSecID = ""
		} else {
			// Open this one, which automatically collapses any other 
			// because ExpandedSecID can only hold one value.
			n.ExpandedSecID = item.ID
		}
	} else {
		// 2. If it's a clickable item, mark as selected
		n.ActiveItemID = item.ID
		
		// Optional: Logic to trigger your eager-loading feature here
		// item.IsLoading = true
		// n.startEagerLoading(ctx, item.ID)

		/*
		// Routing Logic:
        // Assume NavItem.ID matches our route (e.g., "/patients")
        ctx.Navigate(item.ID)
		*/ 
    }
	ctx.Update() // Trigger re-render
}

/*
func (n *LeftNavigation) startEagerLoading(ctx app.Context, id string) {
	// Find the item in our slice to toggle its loading state
	for i := range n.Items {
		if n.Items[i].ID == id {
			// Set loading to true and refresh UI
			n.Items[i].IsLoading = true
			n.Update()

			// Start a background routine
			ctx.Async(func() {
				// Simulate a network delay (e.g., fetching patient data)
				time.Sleep(2 * time.Second)

				// Use Dispatch to return to the UI thread
				ctx.Dispatch(func(ctx app.Context) {
					// Update the state safely back on the main thread
					for j := range n.Items {
						if n.Items[j].ID == id {
							n.Items[j].IsLoading = false
							n.ActiveItemID = id
							break
						}
					}
					// Notify the framework to re-render the component
					n.Update()
					
					log.Println("Data fetch complete for:", id)
				})
			})
			break
		}
	}
}
*/
