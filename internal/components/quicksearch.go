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
	
	// Only allow numeric input (for patient IDs)
	if !q.isValidInput(value) {
		// If invalid, keep the old value
		ctx.Update()
		return
	}
	
	// Limit to max length
	if len(value) > q.MaxLength {
		value = value[:q.MaxLength]
	}
	
	q.SearchTerm = value
	q.Error = "" // Clear error when user types
	ctx.Update()
}

func (q *QuickSearch) onKeyPress(ctx app.Context, e app.Event) {
	// Allow Enter key to trigger search
	keyCode := e.Get("keyCode").Int()
	if keyCode == 13 { // Enter key
		e.PreventDefault()
		q.onSearchClick(ctx, e)
	}
}

func (q *QuickSearch) onSearchClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	
	// Validate input
	if err := q.validate(); err != nil {
		q.Error = err.Error()
		ctx.Update()
		return
	}
	
	// Clear any previous error
	q.Error = ""
	q.IsLoading = true
	ctx.Update()
	
	// If we have a search handler, use it
	if q.OnSearch != nil {
		go q.performSearch(ctx)
	} else {
		go q.mockSearch(ctx)
	}
}

func (q *QuickSearch) performSearch(ctx app.Context) {
	defer func() {
		ctx.Async(func() {
			q.IsLoading = false
			ctx.Update()
		})
	}()
	
	// Call the provided search handler
	err := q.OnSearch(ctx, q.SearchTerm)
	if err != nil {
		ctx.Async(func() {
			q.Error = err.Error()
			ctx.Update()
		})
		return
	}
	
	// Success - clear the search term and update navigation
	ctx.Async(func() {
		q.SearchTerm = ""
		q.Error = "Patient found! Loading..."
		
		// Update application state to trigger navigation change
		// This assumes you have a state management system
		q.updateNavigation(ctx)
		
		ctx.Update()
	})
}

func (q *QuickSearch) mockSearch(ctx app.Context) {
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
		q.SearchTerm = ""
		q.Error = fmt.Sprintf("Mock search for patient ID: %s completed", q.SearchTerm)
		ctx.Update()
	})
}

func (q *QuickSearch) validate() error {
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
	q.SearchTerm = term
	ctx.Update()
}

// ClearSearch clears the search term and error
func (q *QuickSearch) ClearSearch(ctx app.Context) {
	q.SearchTerm = ""
	q.Error = ""
	ctx.Update()
}
