package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Usage:
// For API endpoint: &components.Autocomplete{Endpoint: "/api/search"}
// For mock data: &components.Autocomplete{UseMockData: true}
// For demo: components.NewAutocompleteWithMock()
type Autocomplete struct {
	app.Compo

	// Configuration
	MinChars     int           // Default: 2
	Delay        time.Duration // Default: 300ms
	MaxResults   int           // Default: 10
	Highlight    bool          // Default: true
	Endpoint     string        // API URL (e.g., "/api/search")
	UseMockData  bool          // Use mock data instead of API
	MockData     []string      // Custom mock data
	Placeholder  string        // Input placeholder text

	// Internal State
	query      string
	options    []string        // Store simple strings for display
	results    []AutocompleteResult // Store full result objects
	showPicker bool
	debounce   *time.Timer
	isLoading  bool
	errorMsg   string
	selectedID string
}

type AutocompleteResult struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

// AutocompleteResponse represents the API response structure
type AutocompleteResponse struct {
	Status  string               `json:"status"`
	Data    []AutocompleteResult `json:"data"`
	Query   string               `json:"query,omitempty"`
	Count   int                  `json:"count"`
	HasMore bool                 `json:"has_more"`
}

// NewAutocomplete creates a new autocomplete component
func NewAutocomplete(endpoint string) *Autocomplete {
	return &Autocomplete{
		Endpoint:    endpoint,
		MinChars:    2,
		Delay:       300 * time.Millisecond,
		MaxResults:  10,
		Highlight:   true,
		Placeholder: "Type to search...",
	}
}

// NewAutocompleteWithMock creates a new autocomplete with mock data
func NewAutocompleteWithMock() *Autocomplete {
	return &Autocomplete{
		UseMockData: true,
		MinChars:    2,
		Delay:       300 * time.Millisecond,
		MaxResults:  10,
		Highlight:   true,
		Placeholder: "Type to search (mock data)...",
		MockData:    defaultMockData(),
	}
}

func (a *Autocomplete) OnInit() {
	// Set defaults
	if a.MinChars == 0 {
		a.MinChars = 2
	}
	if a.Delay == 0 {
		a.Delay = 300 * time.Millisecond
	}
	if a.MaxResults == 0 {
		a.MaxResults = 10
	}
	if a.Placeholder == "" {
		a.Placeholder = "Type to search..."
	}
	
	// If no mock data provided but using mock mode, use default
	if a.UseMockData && len(a.MockData) == 0 {
		a.MockData = defaultMockData()
	}
}

func (a *Autocomplete) OnInput(ctx app.Context, e app.Event) {
	a.query = ctx.JSSrc().Get("value").String()
	a.errorMsg = "" // Clear any previous error

	if len(a.query) < a.MinChars {
		a.options = nil
		a.results = nil
		a.showPicker = false
		a.isLoading = false
		ctx.Update()
		return
	}

	// Debounce Logic: Reset timer on every keystroke
	if a.debounce != nil {
		a.debounce.Stop()
	}

	a.isLoading = true
	ctx.Update()
	
	a.debounce = time.AfterFunc(a.Delay, func() {
		if a.UseMockData || a.Endpoint == "" {
			a.fetchMockResults(ctx)
		} else {
			a.fetchAPIResults(ctx)
		}
	})
}

func (a *Autocomplete) fetchAPIResults(ctx app.Context) {
	// Build URL with query parameters
	url := fmt.Sprintf("%s?q=%s&limit=%d", a.Endpoint, a.query, a.MaxResults)
	
	ctx.Async(func() {
		res, err := http.Get(url)
		if err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				a.errorMsg = fmt.Sprintf("Error fetching results: %v", err)
				a.isLoading = false
				ctx.Update()
			})
			return
		}
		defer res.Body.Close()

		// Read response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				a.errorMsg = fmt.Sprintf("Error reading response: %v", err)
				a.isLoading = false
				ctx.Update()
			})
			return
		}

		// Parse JSON response - try both formats
		// First try the structured response
		var structuredResponse AutocompleteResponse
		if err := json.Unmarshal(body, &structuredResponse); err == nil && structuredResponse.Status == "success" {
			ctx.Dispatch(func(ctx app.Context) {
				a.results = structuredResponse.Data
				// Convert to simple strings for display
				a.options = make([]string, len(a.results))
				for i, result := range a.results {
					a.options[i] = result.Label
				}
				a.showPicker = len(a.results) > 0
				a.isLoading = false
				ctx.Update()
			})
			return
		}

		// Try simple string array response
		var simpleResponse map[string]interface{}
		if err := json.Unmarshal(body, &simpleResponse); err == nil {
			if data, ok := simpleResponse["data"].([]interface{}); ok {
				a.options = make([]string, len(data))
				for i, item := range data {
					if str, ok := item.(string); ok {
						a.options[i] = str
					}
				}
				ctx.Dispatch(func(ctx app.Context) {
					a.showPicker = len(a.options) > 0
					a.isLoading = false
					ctx.Update()
				})
				return
			}
		}

		// Last resort: try as simple string array
		var stringArray []string
		if err := json.Unmarshal(body, &stringArray); err == nil {
			ctx.Dispatch(func(ctx app.Context) {
				a.options = stringArray
				a.showPicker = len(a.options) > 0
				a.isLoading = false
				ctx.Update()
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			a.errorMsg = "Invalid response format from server"
			a.isLoading = false
			ctx.Update()
		})
	})
}

func (a *Autocomplete) fetchMockResults(ctx app.Context) {
	// Simulate API delay
	time.Sleep(100 * time.Millisecond)
	
	// Generate mock results based on query
	queryLower := strings.ToLower(a.query)
	var options []string
	var results []AutocompleteResult
	
	// Base mock data
	mockItems := []struct {
		id    string
		label string
		value string
	}{
		{"1", "aaaa", "value_aaaa"},
		{"2", "aaaaa", "value_aaaaa"},
		{"3", "aaaaaa", "value_aaaaaa"},
		{"4", "bbbb", "value_bbbb"},
		{"5", "bbbbb", "value_bbbbb"},
		{"6", "bbbbbb", "value_bbbbbb"},
		{"7", "Apple", "fruit_apple"},
		{"8", "Banana", "fruit_banana"},
		{"9", "Cherry", "fruit_cherry"},
		{"10", "Date", "fruit_date"},
	}
	
	// Filter based on query
	for _, item := range mockItems {
		if queryLower == "" || strings.Contains(strings.ToLower(item.label), queryLower) {
			options = append(options, item.label)
			results = append(results, AutocompleteResult{
				ID:    item.id,
				Label: item.label,
				Value: item.value,
			})
			if len(options) >= a.MaxResults {
				break
			}
		}
	}
	
	ctx.Dispatch(func(ctx app.Context) {
		a.options = options
		a.results = results
		a.showPicker = len(a.options) > 0
		a.isLoading = false
		ctx.Update()
	})
}

func (a *Autocomplete) onSelect(ctx app.Context, index int) {
	if index < 0 || index >= len(a.options) {
		return
	}
	
	a.query = a.options[index]
	a.showPicker = false
	a.isLoading = false
	a.errorMsg = ""
	
	// Store selected ID if we have results
	if index < len(a.results) {
		a.selectedID = a.results[index].ID
		
		// In go-app v10, we can use Dispatch to send an action
		ctx.Dispatch(func(ctx app.Context) {
			// You can handle the selection here or in parent component
			// For example, trigger navigation or update state
			app.Logf("Autocomplete selected: %s (ID: %s)", a.results[index].Label, a.results[index].ID)
		})
	}
	
	ctx.Update()
}

func (a *Autocomplete) Render() app.UI {
	return app.Div().Class("autocomplete-wrapper").Body(
		// Input field
		app.Div().Class("autocomplete-input-wrapper").Body(
			app.Input().
				Type("text").
				Class("autocomplete-input").
				Value(a.query).
				Placeholder(a.Placeholder).
				OnInput(a.OnInput),
			app.If(a.isLoading,
				func() app.UI {
					return app.Div().Class("autocomplete-spinner").Text("⌛")
				},
			),
		),
		
		// Error message
		app.If(a.errorMsg != "",
			func() app.UI {
				return app.Div().Class("autocomplete-error").Text(a.errorMsg)
			},
		),
		
		// Results picker
		app.If(a.showPicker && len(a.options) > 0,
			func() app.UI {
				return app.Div().Class("autocomplete-picker").Body(
					app.Ul().Class("autocomplete-list").Body(
						app.Range(a.options).Slice(func(i int) app.UI {
							opt := a.options[i]
							return app.Li().
								Class("autocomplete-item").
								OnClick(func(ctx app.Context, e app.Event) {
									a.onSelect(ctx, i)
								}).
								Body(
									a.renderOption(opt),
								)
						}),
					),
				)
			},
		),
		
		// No results message
		app.If(a.showPicker && len(a.options) == 0 && !a.isLoading && a.query != "",
			func() app.UI {
				return app.Div().Class("autocomplete-no-results").Text("No results found")
			},
		),
	)
}

func (a *Autocomplete) renderOption(text string) app.UI {
	if !a.Highlight || a.query == "" {
		return app.Text(text)
	}

	// Split the text to highlight the matching part
	index := strings.Index(strings.ToLower(text), strings.ToLower(a.query))
	if index == -1 {
		return app.Text(text)
	}

	return app.Span().Body(
		app.Text(text[:index]),
		app.Strong().Class("autocomplete-highlight").Text(text[index:index+len(a.query)]),
		app.Text(text[index+len(a.query):]),
	)
}

// SetQuery allows programmatically setting the query
func (a *Autocomplete) SetQuery(query string) {
	a.query = query
	a.Defer(func(ctx app.Context) {
		ctx.Update()
	})
}

// Clear clears the query and results
func (a *Autocomplete) Clear() {
	a.query = ""
	a.options = nil
	a.results = nil
	a.showPicker = false
	a.isLoading = false
	a.errorMsg = ""
	a.selectedID = ""
	a.Defer(func(ctx app.Context) {
		ctx.Update()
	})
}

// GetQuery returns the current query
func (a *Autocomplete) GetQuery() string {
	return a.query
}

// GetSelected returns the selected value (if any)
func (a *Autocomplete) GetSelected() string {
	return a.query
}

// GetSelectedID returns the selected ID (if any)
func (a *Autocomplete) GetSelectedID() string {
	return a.selectedID
}

// defaultMockData returns default mock data for demo purposes
func defaultMockData() []string {
	return []string{
		"Apple",
		"Banana",
		"Cherry",
		"Date",
		"Elderberry",
		"Fig",
		"Grape",
		"Honeydew",
		"Kiwi",
		"Lemon",
		"Mango",
		"Nectarine",
		"Orange",
		"Papaya",
		// Medical terms for clinical app
		"Acetaminophen",
		"Amoxicillin",
		"Aspirin",
		"Ibuprofen",
		"Penicillin",
	}
}
