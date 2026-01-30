// internal/pages/home_test.go
package pages

import (
    "testing"
    "strings"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

func TestHomePage(t *testing.T) {
    // Create the component
    hp := &HomePage{}
    
    // Use app.Test to mount the component
    ui := hp.Render()
    
    // Get the HTML string representation
    html := app.Test(ui)
    
    // Check for expected content
    if !strings.Contains(html, "Welcome to the Prototype") {
        t.Errorf("HomePage did not render expected text. Got:\n%s", html)
    }
    
    // You can also check for specific HTML tags
    if !strings.Contains(html, "<h1>") {
        t.Error("HomePage did not render H1 tag")
    }
}
