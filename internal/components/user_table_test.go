package components

import (
    "testing"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

func TestUserTableRender(t *testing.T) {
    ut := &UserTable{
        Users: []User{
            {Name: "Test1 User", Email: "test1@example.com"},
			{Name: "Test2 User", Email: "test2@example.com"},
        },
        CurrentPage: 1,
        PageSize:    10,
    }
    
    // Create a test dispatcher
    disp := app.NewTestDispatcher(ut)
    
    // Check that the component renders something
    rendered := disp.String()
    if rendered == "" {
        t.Error("UserTable did not render anything")
    }
    
    // You can add more specific checks
    // Note: The actual HTML might be complex, so check for key elements
}
