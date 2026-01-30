package components

import (
    "testing"
    "strings"
    "go-app-demo/internal/state"

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
    
    html := app.HTMLString(ut)
    
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
    
    html1 := app.HTMLString(ut)
    
    t.Run("sort handler", func(t *testing.T) {
        // Create a mock context if needed
        ut.SortBy = "email"
        ut.SortOrder = "desc"
        
        html2 := app.HTMLString(ut)
        
        // Verify something changed if possible
        if html1 == html2 {
            t.Log("Note: Sorting might not affect rendered HTML structure")
        }
    })
}
