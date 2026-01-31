package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
    "go-app-demo/internal/models"
)

type Login struct {
    app.Compo
    UserLoginForm models.UserLoginForm // Your data struct
}

func (l *Login) Render() app.UI {
    return app.Form().OnSubmit(l.Submit).Body(
        &components.FormInput{
            Label:     "Username",
            Value:     l.UserLoginForm.Username,
            OnChanged: func(v string) { l.UserLoginForm.Username = v },
        },
        &components.FormInput{
            Label:     "Password",
            Type:      "password",
            Value:     l.UserLoginForm.Password,
            OnChanged: func(v string) { l.UserLoginForm.Password = v },
        },
        app.Button().Type("submit").Text("Login"),
    )
}

func (l *Login) Submit(ctx app.Context, e app.Event) {
    e.PreventDefault() // Prevent the page from reloading
    
    // Log the data to the browser console for testing
    app.Log("Form submitted:", l.UserLoginForm.Username)
    app.Log("Form submitted:", l.UserLoginForm.Password)

    // Example: Navigate to the home page after "login"
    ctx.Navigate("/home")
}
