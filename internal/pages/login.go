package pages

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"go-app-demo/internal/state"
)

type LoginPage struct {
    app.Compo
    UserLoginForm state.UserLoginForm // Your data struct
}

func (p *LoginPage) Render() app.UI {
    return app.Form().OnSubmit(p.Submit).Body(
        &FormInput{
            Label:     "Username",
            Value:     p.UserLoginForm.Username,
            OnChanged: func(v string) { p.UserLoginForm.Username = v },
        },
        &FormInput{
            Label:     "Password",
            Type:      "password",
            Value:     p.UserLoginForm.Password,
            OnChanged: func(v string) { p.UserLoginForm.Password = v },
        },
        app.Button().Type("submit").Text("Login"),
    )
}
