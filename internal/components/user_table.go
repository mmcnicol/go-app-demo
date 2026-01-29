package components

import (
    "encoding/json"
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/state"
    "net/http"
)

type UserTable struct {
    app.Compo
    
    // Data and Metadata
    Users       []state.User
    TotalRows   int
    
    // Pagination/Sorting State
    CurrentPage int
    PageSize    int
    SortBy      string
    SortOrder   string // "asc" or "desc"
    Loading     bool
}

// Add this constructor function
func NewUserTable(users []state.User) *UserTable {
    return &UserTable{
        Users:       users,
        CurrentPage: 1,
        PageSize:    10,
        SortBy:      "name",
        SortOrder:   "asc",
    }
}

func (t *UserTable) Render() app.UI {
    return app.Div().Body(
        app.Table().Class("table").Body(
            app.THead().Body(
                app.Tr().Body(
                    // Sortable Header
                    app.Th().Text("Name").OnClick(func(ctx app.Context, e app.Event) {
                        t.handleSort(ctx, "name")
                    }),
                    app.Th().Text("Email").OnClick(func(ctx app.Context, e app.Event) {
                        t.handleSort(ctx, "email")
                    }),
                ),
            ),
            app.TBody().Body(
                //app.If(t.Loading, app.Tr().Body(app.Td().ColSpan(2).Text("Loading..."))),
                app.If(t.Loading, func() app.UI {
                    return app.Tr().Body(
                        app.Td().ColSpan(2).Text("Loading..."),
                    )
                }),
                app.Range(t.Users).Slice(func(i int) app.UI {
                    return app.Tr().Body(
                        app.Td().Text(t.Users[i].Name),
                        app.Td().Text(t.Users[i].Email),
                    )
                }),
            ),
        ),
        // Pagination Controls
        app.Div().Class("pagination").Body(
            app.Button().Text("Prev").Disabled(t.CurrentPage <= 1).OnClick(t.onPrev),
            app.Span().Text(fmt.Sprintf(" Page %d ", t.CurrentPage)),
            app.Button().Text("Next").Disabled(len(t.Users) < t.PageSize).OnClick(t.onNext),
        ),
    )
}

func (t *UserTable) handleSort(ctx app.Context, field string) {
    if t.SortBy == field {
        // Toggle order if clicking the same column
        if t.SortOrder == "asc" { t.SortOrder = "desc" } else { t.SortOrder = "asc" }
    } else {
        t.SortBy = field
        t.SortOrder = "asc"
    }
    t.fetchData(ctx)
}

func (t *UserTable) onNext(ctx app.Context, e app.Event) {
    t.CurrentPage++
    t.fetchData(ctx)
}

func (t *UserTable) onPrev(ctx app.Context, e app.Event) {
    t.CurrentPage--
    t.fetchData(ctx)
}

func (t *UserTable) fetchData(ctx app.Context) {
    t.Loading = true
    ctx.Update()

    ctx.Async(func() {
        // Build URL with query params
        url := fmt.Sprintf("/api/users?page=%d&size=%d&sort=%s&order=%s",
            t.CurrentPage, t.PageSize, t.SortBy, t.SortOrder)

        res, err := http.Get(url)
        if err != nil {
            return
        }
        defer res.Body.Close()

        var result struct {
            Data  []state.User `json:"data"`
            Total int          `json:"total"`
        }
        json.NewDecoder(res.Body).Decode(&result)

        // Sync back to UI thread
        ctx.Dispatch(func(ctx app.Context) {
            t.Users = result.Data
            t.TotalRows = result.Total
            t.Loading = false
        })
    })
}
