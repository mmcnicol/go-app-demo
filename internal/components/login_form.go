// login_form.go
package components

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
)

// LoginHandler defines the signature for login callback functions
type LoginHandler func(ctx app.Context, username string, password string) (*models.User, error)

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
				Body(
					app.If(l.IsLoading,
						func() app.UI {
							return app.Text("Logging in...")
						},
					).Else(
						func() app.UI {
							return app.Text("Login")
						},
					),
				).
				Disabled(l.IsLoading),
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

	// Validate inputs
	err := l.validate()
	if err != nil {
		l.Error = err.Error()
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
	// Call the provided login handler
	user, err := l.OnLogin(ctx, l.Username, l.Password)

	// Use Async to safely update UI from goroutine
	ctx.Async(func() {
		l.IsLoading = false

		if err != nil {
			l.Error = err.Error()
		} else if user != nil {
			l.Error = "Login successful! Redirecting..."
			// Clear form
			l.Username = ""
			l.Password = ""
			// Note: The parent component will handle navigation
		}

		ctx.Update()
	})
}

func (l *LoginForm) mockLogin(ctx app.Context) {
	// Simulate network delay
	time.Sleep(1 * time.Second)

	// Use Async to safely update UI from goroutine
	ctx.Async(func() {
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

func (l *LoginForm) validate() error {
	if l.Username == "" {
		return fmt.Errorf("username is required")
	}
	if l.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(l.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

/*
// Alternative version using app.Handle for cleaner error handling
func (l *LoginForm) performLoginV2(ctx app.Context) {
	// Call the provided login handler
	user, err := l.OnLogin(ctx, l.Username, l.Password)

	// Handle the result
	ctx.Handle(func() {
		l.IsLoading = false

		if err != nil {
			l.Error = err.Error()
			ctx.Update()
			return
		}

		if user != nil {
			l.Error = "Login successful!"
			l.Username = ""
			l.Password = ""
			ctx.Update()

			// Optionally navigate after a delay
			go func() {
				time.Sleep(1 * time.Second)
				ctx.Navigate("/") // Navigate to home
			}()
		}
	})
}
*/

/*
// For v10, you can also use Defer with the context directly in some cases
func (l *LoginForm) performLoginV3(ctx app.Context) {
	// Defer is available on context but usually for cleanup
	defer ctx.Defer(func() {
		app.Log("Login attempt completed")
	})

	_, err := l.OnLogin(ctx, l.Username, l.Password)

	// Use Dispatch for immediate UI updates
	ctx.Dispatch(func(ctx app.Context) {
		l.IsLoading = false
		
		if err != nil {
			l.Error = err.Error()
		} else {
			l.Error = "Login successful!"
			l.Username = ""
			l.Password = ""
		}
	})
}
*/
