package components

import (
	"bytes"
	"encoding/json"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"net/http"
)

type UserForm struct {
	app.Compo

	// Form State
	name         string
	email        string
	isSubmitting bool
	message      string
}

// OnInputChange updates the struct field when the user types
func (f *UserForm) OnInputChange(ctx app.Context, e app.Event) {
	val := ctx.JSSrc().Get("value").String()
	id := ctx.JSSrc().Get("id").String()

	switch id {
	case "name":
		f.name = val
	case "email":
		f.email = val
	}
}

func (f *UserForm) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault() // Stop the browser from reloading the page
	app.Log("[UserForm OnSubmit] ")	
	f.isSubmitting = true
	f.message = "Sending data..."
	ctx.Update()

	// Logic runs in a background goroutine so the UI doesn't freeze
	ctx.Async(func() {
		app.Log("[UserForm OnSubmit] in Async")
		data := map[string]string{
			"name":  f.name,
			"email": f.email,
		}
		
		jsonData, _ := json.Marshal(data)
		
		// Sending the data to your server-side API
		resp, err := http.Post("/api/users", "application/json", bytes.NewBuffer(jsonData))
		
		// Marshalling back to the UI thread
		ctx.Dispatch(func(ctx app.Context) {
			f.isSubmitting = false
			if err != nil || resp.StatusCode != 200 {
				f.message = "Error: Could not save user."
			} else {
				f.message = "User saved successfully!"
				f.name = ""  // Clear form
				f.email = ""
			}
		})
	})
}

func (f *UserForm) Render() app.UI {
	return app.Div().Class("container").Body(
		app.H2().Text("Add New User"),
		app.Form().OnSubmit(f.OnSubmit).Body(
			app.Div().Body(
				app.Label().For("name").Text("Name: "),
				app.Input().ID("name").Value(f.name).OnChange(f.OnInputChange),
			),
			app.Div().Body(
				app.Label().For("email").Text("Email: "),
				app.Input().ID("email").Type("email").Value(f.email).OnChange(f.OnInputChange),
			),
			app.Button().Type("submit").Disabled(f.isSubmitting).Text("Submit"),
		),
		app.If(f.message != "", func() app.UI {
    		return app.P().Text(f.message)
		}),
	)
}
