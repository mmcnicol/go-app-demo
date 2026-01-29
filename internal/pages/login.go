package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
    "go-app-demo/internal/state"
)

type LoginPage struct {
    app.Compo
    UserLoginForm state.UserLoginForm // Your data struct
}

func (p *LoginPage) Render() app.UI {
    return app.Form().OnSubmit(p.Submit).Body(
        &components.FormInput{
            Label:     "Username",
            Value:     p.UserLoginForm.Username,
            OnChanged: func(v string) { p.UserLoginForm.Username = v },
        },
        &components.FormInput{
            Label:     "Password",
            Type:      "password",
            Value:     p.UserLoginForm.Password,
            OnChanged: func(v string) { p.UserLoginForm.Password = v },
        },
        app.Button().Type("submit").Text("Login"),
    )
}

func (p *LoginPage) Submit(ctx app.Context, e app.Event) {
    e.PreventDefault() // Prevent the page from reloading
    
    // Log the data to the browser console for testing
    app.Log("Form submitted:", p.UserLoginForm.Username)

    // Example: Navigate to the home page after "login"
    ctx.Navigate("/")
}
