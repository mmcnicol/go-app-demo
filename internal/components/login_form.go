package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type LoginForm struct {
    app.Compo
    Username string
    Password string
    Error    string
}

func NewLoginForm() *LoginForm {
    return &LoginForm{}
}

func (l *LoginForm) Render() app.UI {
    return app.Div().Class("login-form").Body(
        app.H2().Text("Login"),
        app.If(l.Error != "",
            func() app.UI {
                return app.Div().Class("error").Text(l.Error)
            },
        ),
        app.Form().OnSubmit(l.onSubmit).Body(
            app.Div().Class("form-group").Body(
                app.Label().Text("Username"),
                app.Input().
                    Type("text").
                    Value(l.Username).
                    OnChange(l.onUsernameChange),
            ),
            app.Div().Class("form-group").Body(
                app.Label().Text("Password"),
                app.Input().
                    Type("password").
                    Value(l.Password).
                    OnChange(l.onPasswordChange),
            ),
            app.Button().
                Type("submit").
                Text("Login"),
        ),
    )
}

func (l *LoginForm) onUsernameChange(ctx app.Context, e app.Event) {
    l.Username = ctx.JSSrc().Get("value").String()
    ctx.Update()
}

func (l *LoginForm) onPasswordChange(ctx app.Context, e app.Event) {
    l.Password = ctx.JSSrc().Get("value").String()
    ctx.Update()
}

func (l *LoginForm) onSubmit(ctx app.Context, e app.Event) {
    e.PreventDefault()
    app.Log("Form submitted:", l.Username)
    app.Log("Form submitted:", l.Password)

    // You would call the parent page's login method here
    // For now, just show a message
    l.Error = "Login functionality not implemented"

    /*
    // Example: Navigate to the home page after "login"
    ctx.Navigate("/home")
    */

    ctx.Update()
}
