package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
	"go-app-demo/internal/state"
)

type LeftNavigation struct {
	app.Compo
	Items         []models.NavItem
	ActiveItemID  string
	ExpandedSecID string
	isMounted     bool // Managed internally
	ctx app.Context
}

func NewLeftNavigation() *LeftNavigation {
	return &LeftNavigation{}
}

func (n *LeftNavigation) OnMount(ctx app.Context) {
	n.ctx = ctx
	
	ctx.ObserveState(state.NavItemsKey, &n.Items).
		OnChange(func() {
			app.Log("[LeftNavigation] Nav Items changed")
			n.Refresh(ctx)
		})
	ctx.ObserveState(state.ExpandedSectionKey, &n.ExpandedSecID).
		OnChange(func() {
			app.Log("[LeftNavigation] Expanded Section changed")
			n.Refresh(ctx)
		})
	ctx.ObserveState(state.ActiveItemKey, &n.ActiveItemID).
		OnChange(func() {
			app.Log("[LeftNavigation] Active Item changed")
			n.Refresh(ctx)
		})
	n.isMounted = true
	n.Refresh(ctx)
}

// Refresh updates the component
func (n *LeftNavigation) Refresh(ctx app.Context) {
	ctx.Update()
}

func (n *LeftNavigation) Render() app.UI {
	app.Log("[LeftNavigation Render] Called")
	app.Log("[LeftNavigation Render] Called, isMounted:", n.isMounted)
	
	// Don't render anything until mounted
	if !n.isMounted {
		app.Log("[LeftNavigation Render] Not mounted yet, returning loading/empty")
		return app.Div().Class("nav-sidebar loading")
	}

	/*
	var items *[]models.NavItem
    ctx.ObserveState(state.GopherDemographicsObservable).Value(&items)
	*/

	return app.Nav().Class("nav-sidebar").Body(
		app.Ul().Class("nav-list").Body(
			app.Range(n.Items).Slice(func(i int) app.UI {
				return n.renderItem(n.Items[i])
			}),
		),
	)
}

// renderItem creates the UI for an individual item or a section
func (n *LeftNavigation) renderItem(item models.NavItem) app.UI {

	hasChildren := len(item.Children) > 0
	
	/*
	var expandedSecID string
	var activeItemID string
    ctx.ObserveState(state.ExpandedSectionObservable).Value(&expandedSecID)
	ctx.ObserveState(state.ActiveItemObservable).Value(&activeItemID)
	*/

	isExpanded := n.ExpandedSecID == item.ID
	isSelected := n.ActiveItemID == item.ID
	
	// Build conditional class string
	itemClass := "nav-item"
	if isSelected {
		itemClass += " selected"
	}
	if hasChildren {
		itemClass += " has-children"
	}
	if isChild {
		itemClass += " child-item"
	}

	return app.Li().Body(
		app.Div().
			Class(itemClass).
			OnClick(func(ctx app.Context, e app.Event) {
				app.Log("[LeftNavigation handleItemClick] Clicked on item:", item.ID, "Label:", item.Label)
				e.PreventDefault() // Prevent default behavior
				//e.StopPropagation() // Stop event bubbling
				n.handleItemClick(ctx, item)
			}).
			Body(
				// Icon or Spinner (Eager Loading Logic)
				app.Div().Class("nav-icon").Body(
					app.If(item.IsLoading,
						func() app.UI {
							return app.Span().Class("spinner")
						},
					).Else(
						func() app.UI {
							return app.I().Class(item.Icon)
						},
					),
				),
				app.Span().Class("nav-label").Text(item.Label),
				// Chevron indicator for sections
				app.If(hasChildren,
					func() app.UI {
						chevronClass := "fas ml-auto"
						if isExpanded {
							chevronClass += " fa-chevron-down"
						} else {
							chevronClass += " fa-chevron-right"
						}
						app.Log("[LeftNavigation renderItem]   Chevron class:", chevronClass)
						return app.I().Class(chevronClass)
					},
				),
			),
		// Submenu (The Accordion Content)
		app.If(hasChildren && isExpanded,
			func() app.UI {
				app.Log("[LeftNavigation renderItem]   Rendering children for:", item.ID)
				return app.Ul().Class("nav-submenu").Body(
					app.Range(item.Children).Slice(func(j int) app.UI {
						return n.renderItem(item.Children[j])
					}),
				)
			},
		),
	)
}

func (n *LeftNavigation) handleItemClick(ctx app.Context, item models.NavItem) {
	
	//ctx.SetState(state.ExpandedSectionObservable, userSettings.defaultNavSectionGopherContext)
	//ctx.SetState(state.ActiveItemObservable, userSettings.defaultNavItemNonGopherContext)
	
	var expandedSecID string
	var activeItemID string
    ctx.ObserveState(state.ExpandedSectionObservable).Value(&expandedSecID)
	ctx.ObserveState(state.ActiveItemObservable).Value(&activeItemID)
	
	// 1. If it's a section (has children), toggle accordion
	if len(item.Children) > 0 {
		app.Log("[LeftNavigation handleItemClick] Item has children, toggling accordion")
		if expandedSecID == item.ID {
			// If already open, close it
			app.Log("[LeftNavigation handleItemClick] Closing section:", item.ID)
			//n.ExpandedSecID = ""
			ctx.SetState(state.ExpandedSectionObservable, nil)
		} else {
			// Open this one, which automatically collapses any other
			// because ExpandedSecID can only hold one value.
			app.Log("[LeftNavigation handleItemClick] Opening section:", item.ID)
			//expandedSecID = item.ID
			ctx.SetState(state.ExpandedSectionObservable, item.ID)
			//activeItemID = n.GetFirstChildID(item)
			ctx.SetState(state.ActiveItemObservable, n.GetFirstChildID(item))
			
			/*
			if n.ActiveItemID == "" {
				app.Log("[LeftNavigation handleItemClick] n.ActiveItemID is empty string!")
			} else {
				app.Log("[LeftNavigation handleItemClick] n.ActiveItemID set to:", n.ActiveItemID)
			}
			*/
		}
	} else {
		// 2. If it's a clickable item, mark as selected
		app.Log("[LeftNavigation handleItemClick] Setting ActiveItemID to:", item.ID)
		//activeItemID = item.ID
		ctx.SetState(state.ActiveItemObservable, item.ID)

		/*
		// Navigation
		if item.Route != "" {
			app.Log("[LeftNavigation handleItemClick] Navigating to route:", item.Route)
			//ctx.Navigate(item.Route)
		} else {
			app.Log("[LeftNavigation handleItemClick] Navigating to ID:", item.ID)
			//ctx.Navigate(item.ID)
		}
		*/
	}
	
	app.Log("[LeftNavigation handleItemClick] END - New ExpandedSecID:", n.ExpandedSecID, "New ActiveItemID:", n.ActiveItemID)
	
	// IMPORTANT: Use ctx.Async to ensure proper state update
	//ctx.Async(func() {
		//app.Log("[LeftNavigation handleItemClick async] Triggering update prerequisite using async")

		// Update the navigation context
		ctx.Dispatch(func(ctx app.Context) {
			app.Log("[LeftNavigation handleItemClick async] Triggering update using dispatch")
			ctx.Update()
		})
		
	//})
}

// GetFirstChildID returns the ID of the first child NavItem if it exists
func (n *LeftNavigation) GetFirstChildID(item models.NavItem) string {
    if len(item.Children) > 0 {
        return item.Children[0].ID
    }
    return "" // No children
}

/*
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
*/

/*
// SetItems allows updating navigation items dynamically
func (n *LeftNavigation) SetItems(items []models.NavItem) {
	app.Log("[LeftNavigation SetItems] Called with", len(items), "items")
	n.Items = items
	
	// Re-initialize state based on new items
	n.initializeState()
	
	//ctx.Async(func() {
	//	ctx.Update()
	//})
}
*/

/*
// SetActiveItem allows programmatically setting active item
func (n *LeftNavigation) SetActiveItem(ctx app.Context, itemID string) {
	app.Log("[LeftNavigation SetActiveItem] Setting ActiveItemID to:", itemID)
	n.ActiveItemID = itemID
	
	ctx.Async(func() {
		ctx.Update()
	})
}
*/

/*
// SetExpandedSection allows programmatically setting expanded section
func (n *LeftNavigation) SetExpandedSection(ctx app.Context, sectionID string) {
	app.Log("[LeftNavigation SetExpandedSection] Setting ExpandedSecID to:", sectionID)
	n.ExpandedSecID = sectionID
	
	ctx.Async(func() {
		ctx.Update()
	})
}
*/
