package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type FormInput struct {
    app.Compo
    Label       string
    Type        string // text, email, password
    Value       string
    Placeholder string
    // Callback to update the parent
    OnChanged   func(string) 
}

func (i *FormInput) Render() app.UI {
    return app.Div().Class("form-group").Body(
        app.Label().Text(i.Label),
        app.Input().
            Type(i.Type).
            Class("form-control").
            Value(i.Value).
            Placeholder(i.Placeholder).
            OnChange(func(ctx app.Context, e app.Event) {
                val := ctx.JSSrc().Get("value").String()
                i.OnChanged(val) // Send data back to parent
            }),
    )
}
