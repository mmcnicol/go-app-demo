package pages

import (
	"testing"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"net/url" 
	"strings"
)

func TestHomePage(t *testing.T) {
    // Initialize your component
    ut := &HomePage{}

    // Use go-app's internal testing tools to "mount" the component
    // and verify the rendered HTML structure
    disp := app.NewTestDispatcher(ut)
    disp.Nav(&url.URL{Path: "/home"})

    // You can then assert that specific text or tags exist in the UI
    if !strings.Contains(disp.String(), "Welcome to the Prototype") {
        t.Error("HomePage did not render H1")
    }
}
