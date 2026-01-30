package components

import (
    "testing"
    "strings"
    "go-app-demo/internal/state"

    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

func TestGopherBannerRender(t *testing.T) {
	
	gopherDemographics := &state.GopherDemographics{
		GopherId: "1"
        Name: "Test",
        DateOfBirth: time.Date(1949, time.January, 2, 0, 0, 0, 0, time.UTC),
    }

	gopherBanner := NewGopherBanner(gopherDemographics)
    
    html := app.HTMLString(gopherBanner)
    
    // Basic checks
    if html == "" {
        t.Error("GopherBanner did not render anything")
    }
    
    // Check for table structure
    if !strings.Contains(html, "<div") {
        t.Error("GopherBanner did not render div")
    }
    
    // Check for gopher demographics; Name
    if !strings.Contains(html, "Test") {
        t.Errorf("GopherBanner did not render Name. Got:\n%s", html)
    }

	// Check for gopher demographics; DateOfBirth
    if !strings.Contains(html, "02/01/1949") {
        t.Errorf("GopherBanner did not render DateOfBirth. Got:\n%s", html)
    }
}
