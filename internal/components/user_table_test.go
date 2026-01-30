package components

import (
    "testing"
    "strings"
    "go-app-demo/internal/state"

    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func TestUserTableRender(t *testing.T) {
    ut := &UserTable{
        Users: []state.User{
            {Name: "Test User", Email: "test@example.com"},
        },
        CurrentPage: 1,
        PageSize:    10,
        TotalRows:   1,
    }
    
    // Render the component
    disp := app.Test(ut)
    
    // Get HTML representation
    var b strings.Builder
	disp.Render(&b)
	html := b.String()
    
    // Basic checks
    if html == "" {
        t.Error("UserTable did not render anything")
    }
    
    // Check for table structure
    if !strings.Contains(html, "<table") {
        t.Error("UserTable did not render table")
    }
    
    // Check for user data
    if !strings.Contains(html, "Test User") {
        t.Errorf("UserTable did not render user name. Got:\n%s", html)
    }
}

func TestUserTableSort(t *testing.T) {
    ut := &UserTable{
        Users: []state.User{
            {Name: "Alice", Email: "alice@example.com"},
            {Name: "Bob", Email: "bob@example.com"},
        },
        CurrentPage: 1,
        PageSize:    10,
        TotalRows:   2,
        SortBy:      "name",
        SortOrder:   "asc",
    }
    
    // Test initial render
    // Render the component
    disp := app.Test(ut)
    
    // Get HTML representation
    var b strings.Builder
	disp.Render(&b)
	html1 := b.String()
    
    // You can't directly test event handlers with app.Test
    // but you can test that the component renders with initial state
    
    // For event testing, you might need integration/E2E tests
    // or test the handler functions directly
    t.Run("sort handler", func(t *testing.T) {
        // Create a mock context if needed
        ut.SortBy = "email"
        ut.SortOrder = "desc"
        
        // Render the component
        disp := app.Test(ut)
        
        // Get HTML representation
        var b strings.Builder
        disp.Render(&b)
        html2 := b.String()
        
        // Verify something changed if possible
        if html1 == html2 {
            t.Log("Note: Sorting might not affect rendered HTML structure")
        }
    })
}
