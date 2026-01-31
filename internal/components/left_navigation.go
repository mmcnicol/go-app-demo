package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
)

type LeftNavigation struct {
	app.Compo
	Items        []models.NavItem
	ActiveItemID string // Managed internally
	ExpandedSecID string // Managed internally
}

// NewLeftNavigation creates a new LeftNavigation component
func NewLeftNavigation(items []models.NavItem) *LeftNavigation {
	app.Log("[LeftNavigation New] Creating with", len(items), "items")
	
	nav := &LeftNavigation{
		Items: items,
	}
	
	// Initialize state immediately in constructor, not in OnMount
	nav.initializeState()
	
	return nav
}

// initializeState sets initial ActiveItemID and ExpandedSecID based on items
func (n *LeftNavigation) initializeState() {
	app.Log("[LeftNavigation initializeState] Called")
	
	// Set initial ActiveItemID based on the first leaf item (non-parent)
	if len(n.Items) > 0 {
		if firstLeaf := n.findFirstLeafItem(n.Items); firstLeaf != nil {
			n.ActiveItemID = firstLeaf.ID
			app.Log("[LeftNavigation initializeState] Setting initial ActiveItemID:", n.ActiveItemID)
		}
	}
	
	// Find and expand sections with IsDefaultExpanded = true
	for _, item := range n.Items {
		app.Log("[LeftNavigation initializeState] Checking item ID:", item.ID, "IsDefaultExpanded:", item.IsDefaultExpanded)
		if item.IsDefaultExpanded && len(item.Children) > 0 {
			n.ExpandedSecID = item.ID
			app.Log("[LeftNavigation initializeState] Setting ExpandedSecID to:", item.ID, "based on IsDefaultExpanded")
			break // Only expand one section by default
		}
	}
	
	app.Log("[LeftNavigation initializeState] Final state - ActiveItemID:", n.ActiveItemID, "ExpandedSecID:", n.ExpandedSecID)
}

// OnMount handles additional initialization if needed
func (n *LeftNavigation) OnMount(ctx app.Context) {
	app.Log("[LeftNavigation OnMount] Called")
	app.Log("[LeftNavigation OnMount] Current ActiveItemID:", n.ActiveItemID)
	app.Log("[LeftNavigation OnMount] Current ExpandedSecID:", n.ExpandedSecID)
	
	// If state wasn't initialized (shouldn't happen with new constructor), initialize it
	if n.ActiveItemID == "" && n.ExpandedSecID == "" {
		app.Log("[LeftNavigation OnMount] State not initialized, initializing now")
		n.initializeState()
	}
}

// Helper to find first leaf item (non-parent) in the navigation tree
func (n *LeftNavigation) findFirstLeafItem(items []models.NavItem) *models.NavItem {
	for _, item := range items {
		if len(item.Children) == 0 {
			return &item
		}
		// If item has children, recursively check children
		if child := n.findFirstLeafItem(item.Children); child != nil {
			return child
		}
	}
	return nil
}

// OnNav handles navigation changes from browser history
func (n *LeftNavigation) OnNav(ctx app.Context) {
	app.Log("[LeftNavigation OnNav] Called")
	// Automatically highlight the item that matches the current URL
	currPath := ctx.Page().URL().Path

	// Find matching item in the navigation tree
	if n.findAndSetActiveItem(currPath, n.Items) {
		ctx.Update()
	}
}

func (n *LeftNavigation) Render() app.UI {
	app.Log("[LeftNavigation Render] Called")
	app.Log("[LeftNavigation Render] Items count:", len(n.Items))
	app.Log("[LeftNavigation Render] ActiveItemID:", n.ActiveItemID)
	app.Log("[LeftNavigation Render] ExpandedSecID:", n.ExpandedSecID)
	
	// Log all item IDs for debugging
	for i, item := range n.Items {
		app.Log("[LeftNavigation Render] Item", i, "ID:", item.ID, "Label:", item.Label, "IsDefaultExpanded:", item.IsDefaultExpanded)
		if len(item.Children) > 0 {
			app.Log("[LeftNavigation Render]   Children:", len(item.Children))
			for j, child := range item.Children {
				app.Log("[LeftNavigation Render]     Child", j, "ID:", child.ID, "Label:", child.Label)
			}
		}
	}
	
	return app.Nav().Class("nav-sidebar").Body(
		app.Ul().Class("nav-list").Body(
			app.Range(n.Items).Slice(func(i int) app.UI {
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
	
	app.Log("[LeftNavigation renderItem] Rendering item:", item.ID, "Label:", item.Label)
	app.Log("[LeftNavigation renderItem]   hasChildren:", hasChildren)
	app.Log("[LeftNavigation renderItem]   isExpanded:", isExpanded, "(ExpandedSecID:", n.ExpandedSecID, "item.ID:", item.ID, ")")
	app.Log("[LeftNavigation renderItem]   isSelected:", isSelected)
	
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
						return n.renderItem(item.Children[j], true)
					}),
				)
			},
		),
	)
}

func (n *LeftNavigation) handleItemClick(ctx app.Context, item models.NavItem) {
	app.Log("[LeftNavigation handleItemClick] START")
	app.Log("[LeftNavigation handleItemClick] Item ID:", item.ID, "Label:", item.Label)
	app.Log("[LeftNavigation handleItemClick] Current ExpandedSecID:", n.ExpandedSecID)
	app.Log("[LeftNavigation handleItemClick] Current ActiveItemID:", n.ActiveItemID)
	app.Log("[LeftNavigation handleItemClick] Has children:", len(item.Children) > 0)
	
	// 1. If it's a section (has children), toggle accordion
	if len(item.Children) > 0 {
		app.Log("[LeftNavigation handleItemClick] Item has children, toggling accordion")
		if n.ExpandedSecID == item.ID {
			// If already open, close it
			app.Log("[LeftNavigation handleItemClick] Closing section:", item.ID)
			n.ExpandedSecID = ""
		} else {
			// Open this one, which automatically collapses any other
			// because ExpandedSecID can only hold one value.
			app.Log("[LeftNavigation handleItemClick] Opening section:", item.ID)
			n.ExpandedSecID = item.ID
		}
	} else {
		// 2. If it's a clickable item, mark as selected
		app.Log("[LeftNavigation handleItemClick] Setting ActiveItemID to:", item.ID)
		n.ActiveItemID = item.ID

		// Navigation
		if item.Route != "" {
			app.Log("[LeftNavigation handleItemClick] Navigating to route:", item.Route)
			ctx.Navigate(item.Route)
		} else {
			app.Log("[LeftNavigation handleItemClick] Navigating to ID:", item.ID)
			ctx.Navigate(item.ID)
		}
	}
	
	app.Log("[LeftNavigation handleItemClick] END - New ExpandedSecID:", n.ExpandedSecID, "New ActiveItemID:", n.ActiveItemID)
	
	// IMPORTANT: Use ctx.Async to ensure proper state update
	ctx.Async(func() {
		app.Log("[LeftNavigation handleItemClick async] Triggering update")
		ctx.Update()
	})
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

// SetItems allows updating navigation items dynamically
func (n *LeftNavigation) SetItems(ctx app.Context, items []models.NavItem) {
	app.Log("[LeftNavigation SetItems] Called with", len(items), "items")
	n.Items = items
	
	// Re-initialize state based on new items
	n.initializeState()
	
	/*
	ctx.Async(func() {
		ctx.Update()
	})
	*/
}

// SetActiveItem allows programmatically setting active item
func (n *LeftNavigation) SetActiveItem(ctx app.Context, itemID string) {
	app.Log("[LeftNavigation SetActiveItem] Setting ActiveItemID to:", itemID)
	n.ActiveItemID = itemID
	
	ctx.Async(func() {
		ctx.Update()
	})
}

// SetExpandedSection allows programmatically setting expanded section
func (n *LeftNavigation) SetExpandedSection(ctx app.Context, sectionID string) {
	app.Log("[LeftNavigation SetExpandedSection] Setting ExpandedSecID to:", sectionID)
	n.ExpandedSecID = sectionID
	
	ctx.Async(func() {
		ctx.Update()
	})
}
