package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

/*
Key Features:

Input Validation:
* Exactly 10 characters required
* Only numeric input allowed
* Real-time validation feedback

User Experience:
* Loading spinner during search
* Error messages for invalid input
* Enter key support
* Max length enforcement

Integration:
* Callback function for search logic
* Right-aligned in application banner
* Works alongside logout button

Styling:
* Responsive design
* Consistent with clinical app theme
* Clear visual feedback

This implementation provides a professional quick search component suitable for a clinical application.
*/

// QuickSearchHandler defines the signature for quick search callback functions
type QuickSearchHandler func(ctx app.Context, searchTerm string) error

type QuickSearch struct {
	app.Compo
	SearchTerm string
	Error      string
	OnSearch   QuickSearchHandler
	IsLoading  bool
	MaxLength  int
	MinLength  int
}

// NewQuickSearch creates a new quick search component
func NewQuickSearch(onSearch QuickSearchHandler) *QuickSearch {
	return &QuickSearch{
		OnSearch:  onSearch,
		MaxLength: 10, // 10-digit patient identifier
		MinLength: 10, // Must be exactly 10 characters
	}
}

func (q *QuickSearch) Render() app.UI {
	return app.Div().Class("quick-search").Body(
		app.Div().Class("search-container").Body(
			app.Input().
				Type("text").
				Class("search-input").
				Value(q.SearchTerm).
				OnChange(q.onSearchTermChange).
				OnKeyPress(q.onKeyPress).
				Placeholder("Enter 10-digit ID").
				MaxLength(q.MaxLength).
				Disabled(q.IsLoading),
			app.Button().
				Class("search-button").
				OnClick(q.onSearchClick).
				Disabled(q.IsLoading).
				Body(
					app.If(q.IsLoading,
						func() app.UI {
							return app.Span().Class("spinner")
						},
					).Else(
						func() app.UI {
							return app.Text("Search")
						},
					),
				),
		),
		app.If(q.Error != "",
			func() app.UI {
				return app.Div().Class("search-error").Text(q.Error)
			},
		),
	)
}

func (q *QuickSearch) onSearchTermChange(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()
	app.Log("[QuickSearch] onSearchTermChange called, value:", value)
	
	// Only allow numeric input (for patient IDs)
	if !q.isValidInput(value) {
		app.Log("[QuickSearch] Invalid input - not numeric")
		// If invalid, keep the old value
		ctx.Update()
		return
	}
	
	// Limit to max length
	if len(value) > q.MaxLength {
		app.Log("[QuickSearch] Truncating value from", len(value), "to", q.MaxLength, "characters")
		value = value[:q.MaxLength]
	}
	
	q.SearchTerm = value
	q.Error = "" // Clear error when user types
	app.Log("[QuickSearch] Search term updated to:", q.SearchTerm)
	ctx.Update()
}

func (q *QuickSearch) onKeyPress(ctx app.Context, e app.Event) {
	keyCode := e.Get("keyCode").Int()
	app.Log("[QuickSearch] onKeyPress called, keyCode:", keyCode)
	// Allow Enter key to trigger search
	if keyCode == 13 { // Enter key
		app.Log("[QuickSearch] Enter key pressed, triggering search")
		e.PreventDefault()
		q.onSearchClick(ctx, e)
	}
}

func (q *QuickSearch) onSearchClick(ctx app.Context, e app.Event) {
	app.Log("[QuickSearch] onSearchClick called")
	e.PreventDefault()
	
	// Validate input
	if err := q.validate(); err != nil {
		app.Log("[QuickSearch] Validation failed:", err.Error())
		q.Error = err.Error()
		ctx.Update()
		return
	}
	
	app.Log("[QuickSearch] Validation passed, search term:", q.SearchTerm)
	
	// Clear any previous error
	q.Error = ""
	q.IsLoading = true
	ctx.Update()
	
	// If we have a search handler, use it
	if q.OnSearch != nil {
		app.Log("[QuickSearch] Starting performSearch")
		go q.performSearch(ctx)
	} else {
		app.Log("[QuickSearch] Starting mockSearch (no handler provided)")
		go q.mockSearch(ctx)
	}
}

func (q *QuickSearch) performSearch(ctx app.Context) {
	app.Log("[QuickSearch] performSearch started")
	defer func() {
		app.Log("[QuickSearch] performSearch completed, setting isLoading=false")
		ctx.Async(func() {
			q.IsLoading = false
			ctx.Update()
		})
	}()
	
	// Call the provided search handler
	app.Log("[QuickSearch] Calling OnSearch handler")
	err := q.OnSearch(ctx, q.SearchTerm)
	if err != nil {
		app.Log("[QuickSearch] OnSearch returned error:", err.Error())
		ctx.Async(func() {
			q.Error = err.Error()
			ctx.Update()
		})
		return
	}
	
	app.Log("[QuickSearch] Search successful, updating UI")
	
	// Success - clear the search term
	ctx.Async(func() {
		q.SearchTerm = ""
		q.Error = "Patient found! Loading..."
		app.Log("[QuickSearch] UI updated: search term cleared, loading message set")
		ctx.Update()
	})
}

func (q *QuickSearch) mockSearch(ctx app.Context) {
	app.Log("[QuickSearch] mockSearch started")
	
	// Simulate network delay
	// In a real app, this would be an API call
	
	ctx.Async(func() {
		q.IsLoading = true
		ctx.Update()
	})
	
	// Simulate API call
	// time.Sleep(500 * time.Millisecond)
	
	ctx.Async(func() {
		q.IsLoading = false
		originalTerm := q.SearchTerm // Save before clearing
		q.SearchTerm = ""
		q.Error = fmt.Sprintf("Mock search for patient ID: %s completed", originalTerm)
		app.Log("[QuickSearch] mockSearch completed, message:", q.Error)
		ctx.Update()
	})
}

func (q *QuickSearch) validate() error {
	app.Log("[QuickSearch] validate called, length:", len(q.SearchTerm), "min:", q.MinLength, "max:", q.MaxLength)
	
	if len(q.SearchTerm) < q.MinLength {
		return fmt.Errorf("Patient ID must be exactly %d digits", q.MinLength)
	}
	
	if len(q.SearchTerm) > q.MaxLength {
		return fmt.Errorf("Patient ID cannot exceed %d digits", q.MaxLength)
	}
	
	// Validate it's numeric
	if !q.isValidInput(q.SearchTerm) {
		return fmt.Errorf("Patient ID must contain only numbers")
	}
	
	app.Log("[QuickSearch] validate passed")
	return nil
}

func (q *QuickSearch) isValidInput(input string) bool {
	// Check if input contains only digits
	for _, char := range input {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// SetSearchTerm allows programmatically setting the search term
func (q *QuickSearch) SetSearchTerm(ctx app.Context, term string) {
	app.Log("[QuickSearch] SetSearchTerm called:", term)
	q.SearchTerm = term
	ctx.Update()
}

// ClearSearch clears the search term and error
func (q *QuickSearch) ClearSearch(ctx app.Context) {
	app.Log("[QuickSearch] ClearSearch called")
	q.SearchTerm = ""
	q.Error = ""
	ctx.Update()
}
