# go-app-demo

A Type-Safe Fullstack Development where both the UI and the API logic share the same Go data structures.

Single language (Go) for frontend and backend reduces "context switching."

Shared Structs: since both the server and the Wasm client are Go, you can share your User or Data structs in a sub-package, ensuring the frontend and backend are always in sync.

Build Speed: Go compiles to Wasm incredibly fast. A typical recompile and refresh often takes less than 2–3 seconds.

Live Reloading: The tool detects a file change, recompiles the app, and refreshes the entire browser page.
* The most popular way to develop with go-app is using [Air](https://github.com/air-verse/air). It watches your .go files, recompiles the Wasm and the server, and restarts the app automatically.

The go-app library has a built-in versioning mechanism. In your app.Handler, you can set a Version string.
* When the Wasm binary is recompiled, you (or your build script) increment this version.
* Your CI should increment a version string in your app.Handler.
* The running PWA detects the version change and can automatically refresh the page to load the new binary.

The go-app library is opinionated: it expects a directory named web in the root. When you run your server, go-app will automatically look here for your .wasm binary, icons, and CSS.

go-app generates a PWA (Progressive Web App), it is perfect for internal tools. Employees can "Install" the tool from the browser onto their desktop or mobile device. It works offline (if you cache the data).

## Project Directory Structure

Building a go-app project requires a slightly different directory structure than a standard Go web app because you are essentially managing two builds (the Wasm frontend and the Go server) within a single codebase.

### Suggested Directory Structure

```Plaintext
.
├── cmd/
│   ├── server/          # The entry point for the backend (app.Handler)
│   │   └── main.go
│   └── wasm/            # The entry point for the frontend (WebAssembly)
│       └── main.go      # Calls app.RunWhenOnBrowser()
├── internal/            # Private application logic
│   ├── components/      # Your reusable UI components (Table, Form, etc.)
│   ├── pages/           # High-level page components (Layouts, Routes)
│   └── state/           # Shared data structures and business logic
├── web/                 # Static assets (must be named "web" for go-app)
│   ├── css/             # Custom stylesheets (Bootstrap/Tailwind)
│   └── images/          # Logo, icons, etc.
├── .air.toml            # (Optional) Config for live reloading
├── go.mod
└── Makefile             # Essential for managing the two-step build process
```

## Testing

### Run native Go tests (API/Logic)

```Go
func TestLoginFormHandler(t *testing.T) {
    // ...
}
```

### Run Wasm-targeted unit tests

Since go-app components are standard Go structs, you can write unit tests for them using the testing package. However, because they are designed to run in a browser, you need to tell the Go compiler to target Wasm during the test.

```Go
func TestUserTable(t *testing.T) {
    // Initialize your component
    ut := &userTable{
        users: []User{{ID: 1, Name: "Alice"}},
    }

    // Use go-app's internal testing tools to "mount" the component
    // and verify the rendered HTML structure
    disp := app.NewTestDispatcher(ut)
    disp.Nav(&url.URL{Path: "/"})

    // You can then assert that specific text or tags exist in the UI
    if !strings.Contains(disp.String(), "Alice") {
        t.Error("Table did not render user name")
    }
}
```

### Run end-to-end (E2E) browser UI tests

```Go
func TestCheckoutFlow(t *testing.T) {
    // ...
}
```

## Continuous Integration

Unlike a standard Go binary, you have to build two things:
* The Wasm Binary: This is your frontend.
* The Server Binary: This is the Go executable that serves the Wasm file and handles API requests.

A Typical CI Workflow (GitHub Actions / GitLab CI)
* Lint: Run golangci-lint.
* Test: Run native Go tests (API/Logic).
* Build Frontend: Run GOOS=js GOARCH=wasm go build -o web/app.wasm.
* Test UI: Run Wasm-targeted unit tests.
* Build Server: Build the main Go entry point.
* E2E: Start the app in the background and run browser UI tests (Playwright/Selenium/Cypress/chromedp).

You will also need a Wasm execution engine like node or wasmer installed on your CI runner.

## Security (CSRF & XSS)

This is where go-app (and WebAssembly in general) has a unique advantage.
* Cross-Site Scripting (XSS): go-app is inherently resistant to XSS. When you use .Text(userInput), the library treats that string as plain text, not HTML. It doesn't use innerHTML under the hood, so malicious <script> tags provided by users are never executed.
* CSRF (Cross-Site Request Forgery): 
  * The Frontend side: Since you are writing the HTTP client calls in Go, you have full control over headers. You can easily pull a CSRF token from a cookie or a meta tag and inject it into every http.Request.
  * The Backend side: The go-app server-side handler is a standard http.Handler. You can wrap it in any standard Go security middleware (like gorilla/csrf) just like a traditional web app.
* The Wasm "Sandbox": Your application logic is compiled to a binary. Unlike JavaScript, where a user can easily "Read" your source code in the browser tools to find vulnerabilities, WebAssembly is a compiled format, which adds a layer of "security through obscurity" (though not a primary defense).

## Themes and CSS

Because go-app renders standard HTML5, it is 100% compatible with any CSS framework.
* Tailwind / Bootstrap: This is the most common path. You include the Bootstrap or Tailwind CSS file in your Handler configuration, and then use .Class() in your Go code.
* Scoped Styling: go-app doesn't have "CSS-in-JS" built-in, but you can use the .Style() method for dynamic logic.
* The "Nice" Factor: Since you aren't fighting a framework's default theme, your app looks exactly as good as your CSS skills (or your CSS framework of choice) allow.

An example component using go-app's built-in style system:
```Go
type Button struct {
    app.Compo
    variant string // "primary", "secondary"
}

func (b *Button) Render() app.UI {
    styles := map[string]string{
        "padding": "12px 24px",
        "borderRadius": "4px",
        "fontWeight": "600",
        "cursor": "pointer",
    }
    
    // Conditional styling
    if b.variant == "primary" {
        styles["backgroundColor"] = "#0066CC"
        styles["color"] = "white"
    } else {
        styles["backgroundColor"] = "#F0F0F0"
        styles["color"] = "#333"
    }
    
    return app.Button().
        Text("Click Me").
        Styles(styles)
}
```

## UI Component Library

The Reality: go-app does not come with a built-in library of high-level components (like Material UI or Ant Design).
* The Go-App Philosophy: It provides the primitive HTML elements. You are expected to build your own reusable Go structs for buttons, modals, and inputs.
* The Corporate Solution: Most corporations using go-app create a internal "Design System" package. You create a Go package containing your branded components, and then all your internal apps simply import that package. It’s highly consistent but requires more initial setup than "dragging and dropping" an existing library.

Developers in the go-app (and broader Go-Wasm) ecosystem typically achieve a "Storybook" like feature by implementing a parallel route within their own application that serves as a "Native Custom Go Component Gallery."

## The Architecture

go-app is essentially the "Go way" of doing what React and Web Components do, but with a few philosophical shifts that make it feel very different in practice.
* go-app is very similar to React in terms of how it updates.
  * Declarative UI: Like React’s JSX, you describe what the UI should look like based on the current state. You don't say "find this button and change its color"; you say "if the state is X, the button is red."
  * The Reconciliation Loop: When you call Update(), go-app performs a "diff" between the current UI and the new UI. It only touches the real browser DOM where differences exist. This is the exact same logic as the React Virtual DOM.
  * State Management: Both rely on "Components" as the source of truth for data.

The Difference: In React, you use JavaScript/JSX. In go-app, you use Strongly Typed Go. This means you get compile-time errors for your UI logic, which is a massive win for maintainability.

The biggest thing that sets go-app apart from React is that it is a Progressive Web App (PWA) framework by default.

When you build a go-app project, it automatically generates:
* A Service Worker: To cache your app for offline use.
* A Manifest file: So users can "install" your website as an app on their phone or desktop.
* Server-Side Rendering (SSR): It handles the initial SEO-friendly HTML delivery before the WebAssembly takes over.

In short: It's like React's programming model, but using Go's type safety, compiled into a high-performance binary that works offline like a native mobile app.

## HTML Elements

For every standard browser tag (like <div>, <h1>, <a>, or <input>), go-app provides a corresponding Go function. These functions use method chaining to set attributes, styles, and event handlers.

An example:
```Go
app.Input().Type("text").Placeholder("Name")
```

## Go Custom Components

Unlike standard Web Components (which use the Shadow DOM), go-app components are purely a Go abstraction.

A "Component" in go-app is a Go struct that embeds app.Compo. This struct stores your application state (data) and implements a Render() method to describe how that state looks in HTML.

```Go
type MyCard struct {
    app.Compo
    Title       string
    Description string
}

func (c *MyCard) Render() app.UI {
    // Here we combine standard HTML "blocks" into a custom component
    return app.Div().Class("card").Body(
        app.H2().Text(c.Title),
        app.P().Text(c.Description),
    )
}
```

### Automatic Updates

If you modify the t.users slice and call t.Update(), the library will perform a "diff" and only re-render the rows that changed.

## Conditional Logic

In go-app, conditional logic is handled using standard Go if-else or switch statements directly inside your Render() method.

Because Render() returns an app.UI, you simply use Go logic to decide which component or HTML element to return. The library’s engine will then look at the result and update only the parts of the DOM that changed.

### Basic If-Else Logic

This is the most common way to show or hide elements (e.g., showing a "Loading..." spinner vs. the actual data).

```Go
func (t *userTable) Render() app.UI {
    return app.Div().Body(
        app.H1().Text("User List"),
        
        // Conditional Logic
        app.If(len(t.users) == 0,
            app.P().Text("No users found. Please add some!"),
        ).Else(
            app.Table().Body(
                // ... table rows here
            ),
        ),
    )
}
```

### Using Standard Go Variables

Since Render() is just a function, you can perform complex logic at the top of the function and store the result in a variable.

```Go
func (c *MyComponent) Render() app.UI {
    var display app.UI

    switch c.Status {
    case "success":
        display = app.Text("Operation successful!")
    case "error":
        display = app.Span().Style("color", "red").Text("Something went wrong.")
    default:
        display = app.Text("Pending...")
    }

    return app.Div().Body(
        app.H2().Text("Status Check"),
        display, // Render the variable
    )
}
```

### Short-Circuiting (The "Empty" UI)

If you want to render nothing under a certain condition, you can return nil or an empty app.Div(). However, the cleanest way is often a simple app.If() without an .Else().

```Go
func (p *userPage) Render() app.UI {
    return app.Div().Body(
        app.H1().Text("Dashboard"),
        // Only show the warning if there is an error
        app.If(p.errorMessage != "",
            app.Div().Class("alert-box").Text(p.errorMessage),
        ),
    )
}
```

## Key Tools for Background Tasks

ctx.Async	Starts a goroutine. Use Case: Fetching data from an API, running a timer, or complex calculations.

ctx.Dispatch Runs a function on the UI thread. Use Case: Updating a struct field (state) and triggering a re-render.

ctx.Defer Schedules a function for the next frame. Use Case: Logic that needs to run after the current render is finished.


## Real-World Example: Countdown Timer

```Go
type timerComponent struct {
	app.Compo
	Seconds int
}

func (t *timerComponent) OnMount(ctx app.Context) {
	// Start a background task when the component appears
	ctx.Async(func() {
		for {
			time.Sleep(time.Second)

			// Use Dispatch to safely update the component state
			ctx.Dispatch(func(ctx app.Context) {
				t.Seconds++
				// Update() is automatically called by Dispatch 
				// when the state changes
			})
		}
	})
}

func (t *timerComponent) Render() app.UI {
	return app.H1().Text(fmt.Sprintf("Seconds passed: %d", t.Seconds))
}
```

## Real-World Example: Periodic Data Refresh

Imagine you want your user table to refresh every 30 seconds from a server API.

```Go
func (t *userTable) OnNav(ctx app.Context) {
    ctx.Async(func() {
        for {
            // 1. Get data from server
            newUsers := fetchDataFromServer() 

            // 2. Synchronize back to the UI thread
            ctx.Dispatch(func(ctx app.Context) {
                t.users = newUsers
                // The table automatically re-renders now
            })

            time.Sleep(30 * time.Second)
        }
    })
}
```

## Hyperlinks (app.A())

When you use a link in go-app, it usually performs Client-Side Routing rather than a traditional page reload or a simple AJAX call.
* Internal Links: If you click a link to a route defined in your app.Route, the library intercepts the click. It updates the URL in the address bar and swaps out the components on the screen without ever asking the server for a new HTML page.
* External Links: If the link points outside your app, the browser handles it like a normal navigation.

## Button Clicks (app.Button())

A button click does nothing involving the server unless you explicitly tell it to.

When you write an OnClick handler, you are executing a Go function inside the browser's Wasm runtime. If you need data from a database, you must manually perform an HTTP request (which is the Wasm equivalent of an AJAX/Fetch call).

```Go
func (p *userPage) OnClick(ctx app.Context, e app.Event) {
    // This runs in the browser!
    go func() {
        // Standard Go HTTP call -> translated to AJAX/Fetch
        res, err := http.Get("/api/users")
        if err != nil {
            return
        }
        
        // Update local state and trigger UI refresh
        // ... decode JSON ...
        p.Update() 
    }()
}
```

## Lifecycle Coordination

When you nest components, go-app manages the lifecycle automatically:
* OnMount: Called when the component is inserted into the browser DOM.
* OnNav: Called when the page is navigated to.
* Update: If you call p.Update() on the parent, it triggers a re-render check for all nested children.

## Key Strategies for Complex Layouts

### Component (Nesting)

```Go
type userPage struct {
	app.Compo
	// You can store components as fields to keep their state
	nav     *navbar
	sidebar *sidebar
}

// OnInit is a lifecycle method called when the component is created
func (p *userPage) OnInit() {
	p.nav = &navbar{}
	p.sidebar = &sidebar{}
}

func (p *userPage) Render() app.UI {
	return app.Div().Body(
		p.nav, // Nesting the navbar component
		app.Main().Body(
			p.sidebar, // Nesting the sidebar component
			app.Div().Style("margin-left", "210px").Body(
				app.H2().Text("User Management"),
				// You can also nest components inline
				&userTable{
					users: []User{
						{1, "Alice", "alice@gmail.com"},
						{2, "Bob", "bob@gmail.com"},
					},
				},
			),
		),
	)
}
```

### Passing Data Down (Props)

In go-app, "props" are just exported fields on your struct. In the example above, the userPage passed a slice of User data to the userTable struct during initialization.

### Content Injection (Slots)

Sometimes you want a layout component (like a "Modal" or "Card") that can wrap any content. You do this by adding an app.UI field to your struct:

```Go
type wrapper struct {
	app.Compo
	Content app.UI // This acts like a "slot" or "children"
}

func (w *wrapper) Render() app.UI {
	return app.Div().Class("fancy-border").Body(
		app.H3().Text("Wrapper Title"),
		w.Content, // Render whatever was passed in
	)
}

// Usage:
// &wrapper{ Content: app.Text("Hello World") }
```

## the Observer pattern

In go-app, the Observer pattern is a fundamental alternative to callbacks and is often the preferred approach for managing state changes and communication between components.

The Observer pattern in go-app allows components to subscribe to state changes and get notified when those changes occur, rather than passing callbacks down through the component tree.

Key Components:
* Observable: A source that holds state and can be observed
* Observer: Components that subscribe to the observable
* Notification: When the observable state changes, all observers are notified


