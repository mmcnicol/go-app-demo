package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
)

// LoginHandler defines the signature for login callback functions
type LoginHandler func(ctx app.Context, username, password string) (*models.User, error)

type LoginForm struct {
	app.Compo
	Username   string
	Password   string
	Error      string
	OnLogin    LoginHandler // Callback function
	IsLoading  bool
}

// NewLoginForm creates a new login form with optional login handler
func NewLoginForm(onLogin LoginHandler) *LoginForm {
	return &LoginForm{
		OnLogin: onLogin,
	}
}

func (l *LoginForm) Render() app.UI {
	return app.Div().Class("login-form").Body(
		app.H2().Text("Login"),
		app.If(l.Error != "",
			func() app.UI {
				return app.Div().Class("error-message").Text(l.Error)
			},
		),
		app.Form().OnSubmit(l.onSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("username").Text("Username"),
				app.Input().
					ID("username").
					Type("text").
					Value(l.Username).
					OnChange(l.onUsernameChange).
					Placeholder("Enter username").
					Disabled(l.IsLoading),
			),
			app.Div().Class("form-group").Body(
				app.Label().For("password").Text("Password"),
				app.Input().
					ID("password").
					Type("password").
					Value(l.Password).
					OnChange(l.onPasswordChange).
					Placeholder("Enter password").
					Disabled(l.IsLoading),
			),
			app.Button().
				Type("submit").
				Text(app.If(l.IsLoading, "Logging in...").Else("Login")).
				Disabled(l.IsLoading),
		),
	)
}

func (l *LoginForm) onUsernameChange(ctx app.Context, e app.Event) {
	l.Username = ctx.JSSrc().Get("value").String()
	l.Update()
}

func (l *LoginForm) onPasswordChange(ctx app.Context, e app.Event) {
	l.Password = ctx.JSSrc().Get("value").String()
	l.Update()
}

func (l *LoginForm) onSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
    
	// Validate inputs
	if l.Username == "" || l.Password == "" {
		l.Error = "Username and password are required"
        l.IsLoading = false
		ctx.Update()
		return
	}
	
    app.Log("Form submitted:", l.Username)
    app.Log("Form submitted:", l.Password)
	
    // Show loading state
	l.IsLoading = true
	l.Error = ""
	ctx.Update()
	
	// If we have a login handler, use it
	if l.OnLogin != nil {
		go l.performLogin(ctx)
	} else {
		// Fallback to mock login
		go l.mockLogin(ctx)
	}
}

func (l *LoginForm) performLogin(ctx app.Context) {
	defer func() {
		l.Defer(func() {
			l.IsLoading = false
			ctx.Update()
		})
	}()
	
	// Call the provided login handler
	user, err := l.OnLogin(ctx app.Context, l.Username, l.Password)
	if err != nil {
		l.Defer(func() {
			l.Error = err.Error()
			ctx.Update()
		})
		return
	}
	
	// Success - you could trigger a navigation or state update here
	// Typically the parent component would handle what happens after successful login
	l.Defer(func() {
		l.Error = "Login successful! Redirecting..."
		// Clear form
		l.Username = ""
		l.Password = ""
		ctx.Update()
		
		// In a real app, you might navigate to another page
		// app.Window().Get("location").Set("href", "/dashboard")

        /*
        // Example: Navigate to the home page after "login"
        ctx.Navigate("/home")
        */

	})
}

func (l *LoginForm) mockLogin(ctx app.Context) {
	// Simulate network delay
	app.Dispatch(func() {
		l.IsLoading = true
		ctx.Update()
	})
	
	// Simulate API call
	time.Sleep(1 * time.Second)
	
	l.Defer(func() {
		l.IsLoading = false
		
		// Simple mock validation
		if l.Username == "demo" && l.Password == "demo" {
			l.Error = "Login successful (mock)! Username: demo, Password: demo"
			l.Username = ""
			l.Password = ""
		} else {
			l.Error = "Invalid credentials. Try demo/demo"
		}
		ctx.Update()
	})
}
