package components

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Usage:
// &components.Autocomplete{
//     Endpoint: "/api/search?q=",
//     Highlight: true,
//     Delay: 500 * time.Millisecond,
// }
type Autocomplete struct {
	app.Compo

	// Configuration
	MinChars    int           // Default: 3
	Delay       time.Duration // Default: 600ms
	MaxResults  int           // Default: 100
	Highlight   bool          // Default: true
	Endpoint    string        // API URL (e.g., "/api/search?q=")

	// Internal State
	query      string
	options    []string
	showPicker bool
	debounce   *time.Timer
}

func (a *Autocomplete) OnInit() {
	if a.MinChars == 0 { a.MinChars = 3 }
	if a.Delay == 0 { a.Delay = 600 * time.Millisecond }
	if a.MaxResults == 0 { a.MaxResults = 100 }
}

func (a *Autocomplete) OnInput(ctx app.Context, e app.Event) {
	a.query = ctx.JSSrc().Get("value").String()

	if len(a.query) < a.MinChars {
		a.options = nil
		a.showPicker = false
		return
	}

	// Debounce Logic: Reset timer on every keystroke
	if a.debounce != nil {
		a.debounce.Stop()
	}

	a.debounce = time.AfterFunc(a.Delay, func() {
		a.fetchResults(ctx)
	})
}

func (a *Autocomplete) fetchResults(ctx app.Context) {
	// Call the API
	url := fmt.Sprintf("%s%s&limit=%d", a.Endpoint, a.query, a.MaxResults)
	
	ctx.Async(func() {
		res, err := http.Get(url)
		if err != nil {
			return
		}
		defer res.Body.Close()

		var results []string
		json.NewDecoder(res.Body).Decode(&results)

		ctx.Dispatch(func(ctx app.Context) {
			a.options = results
			a.showPicker = len(a.options) > 0
		})
	})
}

func (a *Autocomplete) onSelect(ctx app.Context, val string) {
	a.query = val
	a.showPicker = false
	// You could trigger a custom event here to notify a parent component
}

func (a *Autocomplete) Render() app.UI {
	return app.Div().Class("autocomplete-wrapper").Body(
		app.Input().
			Type("text").
			Value(a.query).
			Placeholder("Type to search...").
			OnInput(a.OnInput),

		app.If(a.showPicker,
			app.Ul().Class("picklist").Body(
				app.Range(a.options).Slice(func(i int) app.UI {
					opt := a.options[i]
					return app.Li().OnClick(func(ctx app.Context, e app.Event) {
						a.onSelect(ctx, opt)
					}).Body(
						a.renderOption(opt),
					)
				}),
			),
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
		app.Strong().Style("color", "blue").Text(text[index:index+len(a.query)]),
		app.Text(text[index+len(a.query):]),
	)
}
