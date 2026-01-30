package pages

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/maxence-charriere/go-app/v10/pkg/app/testing"
)

func TestHomePage(t *testing.T) {
	// 1. Create the component
	hp := &HomePage{}

	// 2. Setup the test handler
	// In v10, app.Test returns a dispatcher that "mounts" the component
	disp := apptest.Test(hp)

	// 3. Get the HTML representation of the component
	// We use strings.Builder to capture the rendered output
	var b strings.Builder
	disp.Render(&b)
	html := b.String()

	// 4. Assertions
	if !strings.Contains(html, "Welcome to the Prototype") {
		t.Errorf("HomePage did not render expected text. Got:\n%s", html)
	}

	if !strings.Contains(html, "<h1>") {
		t.Error("HomePage did not render H1 tag")
	}
}
