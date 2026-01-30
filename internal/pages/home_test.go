package pages

import (
	"strings"
	"testing"
	
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func TestHomePage(t *testing.T) {
    hp := &HomePage{}
    
    html := app.HTMLString(hp)

    if !strings.Contains(html, "<h1>Welcome") {
        t.Errorf("H1 not found in rendered HTML")
    }
}
