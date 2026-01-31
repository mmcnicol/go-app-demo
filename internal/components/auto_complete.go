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
// For API endpoint: &components.Autocomplete{Endpoint: "/api/search?q="}
// For mock data: &components.Autocomplete{UseMockData: true}
// For demo: components.NewAutocompleteWithMock()
type Autocomplete struct {
	app.Compo

	// Configuration
	MinChars     int           // Default: 2
	Delay        time.Duration // Default: 300ms
	MaxResults   int           // Default: 10
	Highlight    bool          // Default: true
	Endpoint     string        // API URL (e.g., "/api/search?q=")
	UseMockData  bool          // Use mock data instead of API
	MockData     []string      // Custom mock data
	Placeholder  string        // Input placeholder text

	// Internal State
	query      string
	options    []string
	showPicker bool
	debounce   *time.Timer
	isLoading  bool
	errorMsg   string
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
		a.showPicker = false
		a.isLoading = false
		return
	}

	// Debounce Logic: Reset timer on every keystroke
	if a.debounce != nil {
		a.debounce.Stop()
	}

	a.isLoading = true
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
			})
			return
		}

		// Parse JSON response
		var response AutocompleteResponse
		if err := json.Unmarshal(body, &response); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				a.errorMsg = fmt.Sprintf("Error parsing JSON: %v", err)
				a.isLoading = false
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			a.results = response.Data
			a.showPicker = len(a.results) > 0
			a.isLoading = false
		})
	})
}

func (a *Autocomplete) fetchMockResults(ctx app.Context) {
	// Simulate API delay
	time.Sleep(100 * time.Millisecond)
	
	// Generate mock JSON response
	mockResponse := a.generateMockResponse()
	
	ctx.Dispatch(func(ctx app.Context) {
		a.results = mockResponse.Data
		a.showPicker = len(a.results) > 0
		a.isLoading = false
	})
}

func (a *Autocomplete) generateMockResponse() AutocompleteResponse {
	// Generate mock data based on query
	var results []AutocompleteResult
	queryLower := strings.ToLower(a.query)
	
	// Base mock data - you can customize this
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
			results = append(results, AutocompleteResult{
				ID:    item.id,
				Label: item.label,
				Value: item.value,
			})
			if len(results) >= a.MaxResults {
				break
			}
		}
	}
	
	return AutocompleteResponse{
		Status:  "success",
		Data:    results,
		Query:   a.query,
		Count:   len(results),
		HasMore: len(results) >= a.MaxResults,
	}
}

func (a *Autocomplete) onSelect(ctx app.Context, val string) {
	
	// Set the display value based on FieldToShow
	switch a.FieldToShow {
	case "id":
		a.query = result.ID
	case "value":
		a.query = result.Value
	default:
		a.query = result.Label
	}

	a.selectedID = result.ID
	a.showPicker = false
	a.isLoading = false
	a.errorMsg = ""

	// Emit event with full result object
	ctx.Emit("autocomplete-select", map[string]interface{}{
		"id":    result.ID,
		"label": result.Label,
		"value": result.Value,
	})
}

func (a *Autocomplete) getDisplayText(result AutocompleteResult) string {
	switch a.FieldToShow {
	case "id":
		return result.ID
	case "value":
		return result.Value
	default:
		return result.Label
	}
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
		app.If(a.showPicker && len(a.results) > 0,
			func() app.UI {
				return app.Div().Class("autocomplete-picker").Body(
					app.Ul().Class("autocomplete-list").Body(
						app.Range(a.results).Slice(func(i int) app.UI {
							opt := getDisplayText(a.results[i])
							return app.Li().
								Class("autocomplete-item").
								OnClick(func(ctx app.Context, e app.Event) {
									a.onSelect(ctx, opt)
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
	a.Update()
}

// Clear clears the query and results
func (a *Autocomplete) Clear() {
	a.query = ""
	a.options = nil
	a.showPicker = false
	a.isLoading = false
	a.errorMsg = ""
	a.Update()
}

// GetQuery returns the current query
func (a *Autocomplete) GetQuery() string {
	return a.query
}

// GetSelected returns the selected value (if any)
func (a *Autocomplete) GetSelected() string {
	return a.query
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

// Add custom CSS styles
func init() {
	app.Handle(func(ctx app.Context, action app.Action) {
		// This ensures CSS is available
	})
}
