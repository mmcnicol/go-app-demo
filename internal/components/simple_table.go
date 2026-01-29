package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// 1. Define your data structure
type User struct {
    ID    int
    Name  string
    Email string
}

// 2. Define your component
type userTable struct {
    app.Compo
    users []User
}

// 3. Implement the Render method
func (t *userTable) Render() app.UI {
    return app.Table().Body(
        // Table Header
        app.THead().Body(
            app.Tr().Body(
                app.Th().Text("ID"),
                app.Th().Text("Name"),
                app.Th().Text("Email"),
            ),
        ),
        // Table Body with Rows
        app.TBody().Body(
            app.Range(t.users).Slice(func(i int) app.UI {
                u := t.users[i]
                return app.Tr().Body(
                    app.Td().Text(u.ID),
                    app.Td().Text(u.Name),
                    app.Td().Text(u.Email),
                )
            }),
        ),
    )
}
