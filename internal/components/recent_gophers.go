package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
	"time"
)

type RecentGophers struct {
    app.Compo
    recentGophers []models.RecentGopherItem
    SortBy      string
    SortOrder   string // "asc" or "desc"
    Loading     bool
}

func NewRecentGophers(recentGophers []models.RecentGopherItem) *RecentGophers {
    return &RecentGophers{
        recentGophers: recentGophers,
    }
}

func (r *RecentGophers) Render() app.UI {
    return app.Table().Body(
        // Table Header
        app.THead().Body(
            app.Tr().Body(
                // Sortable Header with sort indicators
                app.Th().Text("GopherId").OnClick(func(ctx app.Context, e app.Event) {
                    r.handleSort(ctx, "GopherId")
                }).Class(r.headerClass("GopherId")),
                app.Th().Text("Name").OnClick(func(ctx app.Context, e app.Event) {
                    r.handleSort(ctx, "Name")
                }).Class(r.headerClass("Name")),
                app.Th().Text("DateOfBirth").OnClick(func(ctx app.Context, e app.Event) {
                    r.handleSort(ctx, "DateOfBirth")
                }).Class(r.headerClass("DateOfBirth")),
                app.Th().Text("DateLastAccessed").OnClick(func(ctx app.Context, e app.Event) {
                    r.handleSort(ctx, "DateLastAccessed")
                }).Class(r.headerClass("DateLastAccessed")),
            ),
        ),
        // Table Body with Rows
        app.TBody().Body(
            app.If(r.Loading, func() app.UI {
                return app.Tr().Body(
                    app.Td().ColSpan(3).Text("Loading..."),
                )
            }),
            app.Range(r.gopherDemographics).Slice(func(i int) app.UI {
                row := r.gopherDemographics[i]
				dobFormatted := row.DateOfBirth.Format("02/01/2006")
                lastAccessedFormatted := row.DateLastAccessed.Format("02/01/2006 15:04")
                return app.Tr().Body(
                    app.Td().Text(row.GopherId),
                    app.Td().Text(row.Name),
                    app.Td().Text(dobFormatted),
                    app.Td().Text(lastAccessedFormatted),
                )
            }),
        ),
    )
}

func (r *RecentGophers) headerClass(field string) string {
    if r.SortBy != field {
        return "sortable"
    }
    if r.SortOrder == "asc" {
        return "sorted-asc"
    }
    return "sorted-desc"
}

func (r *RecentGophers) handleSort(ctx app.Context, field string) {
    if r.SortBy == field {
        // Toggle order if clicking the same column
        if r.SortOrder == "asc" { r.SortOrder = "desc" } else { r.SortOrder = "asc" }
    } else {
        r.SortBy = field
        r.SortOrder = "asc"
    }
    t.sortData(ctx)
    ctx.Update() // Trigger re-render
}

func (r *RecentGophers) sortData() {
    if r.SortBy == "" {
        return // No sort selected
    }

    sort.Slice(r.recentGophers, func(i, j int) bool {
        a := r.recentGophers[i]
        b := r.recentGophers[j]

        var less bool
        
        switch r.SortBy {
        case "GopherId":
            less = strings.ToLower(a.GopherId) < strings.ToLower(b.GopherId)
        case "Name":
            less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
        case "DateOfBirth":
            less = a.DateOfBirth.Before(b.DateOfBirth)
        case "DateLastAccessed":
            less = a.DateLastAccessed.Before(b.DateLastAccessed)
        default:
            return false
        }
        
        // Reverse order if descending
        if r.SortOrder == "desc" {
            return !less
        }
        return less
    })
}

func (r *RecentGophers) OnMount(ctx app.Context) {
    r.Loading = true
    ctx.Update()

    r.SortBy = "DateLastAccessed"
    r.SortOrder = "desc"
    r.sortData()
    
    r.Loading = false
    ctx.Update()
}

/*
func (r *RecentGophers) onPatientClick(ctx app.Context, patientID string) {
    // Navigates to e.g., /discharge/MRN-8829
    ctx.Navigate("/discharge/" + patientID)
}
*/

// In the Patient List component
func (r *RecentGophers) onSelect(ctx app.Context, selected models.RecentGopherItem) {
    /*
    // Save the object to the browser's persistent state
    // Persistent State acts like a global "Store" within the browser session.
    ctx.SetState("current-patient", selected)
    ctx.Navigate("/discharge")
    */
}
