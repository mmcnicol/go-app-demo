package components

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/models"
    "strings"
    "sort"
)

// Usage:
// &LabResults{ labResults: data }
type LabResults struct {
	app.Compo
	labResults  []models.LabResultItem
	SortBy      string
    SortOrder   string // "asc" or "desc"
    Loading     bool
}

func NewLabResults(labResults []models.LabResultItem) *LabResults {
    return &LabResults{
        labResults: labResults,
    }
}

func (l *LabResults) Render() app.UI {
	return app.Div().Body(
		app.H3().Text("Lab Results"),
		app.Table().Body(
			// Table Header
			app.THead().Body(
				app.Tr().Body(
					// Sortable Header with sort indicators
					app.Th().Text("ID").OnClick(func(ctx app.Context, e app.Event) {
						l.handleSort(ctx, "ID")
					}).Class(l.headerClass("ID")),
					app.Th().Text("ReportDate").OnClick(func(ctx app.Context, e app.Event) {
						l.handleSort(ctx, "ReportDate")
					}).Class(l.headerClass("ReportDate")),
					app.Th().Text("Subject").OnClick(func(ctx app.Context, e app.Event) {
						l.handleSort(ctx, "Subject")
					}).Class(l.headerClass("Subject")),
				),
			),
			// Table Body with Rows
			app.TBody().Body(
				app.If(l.Loading, func() app.UI {
					return app.Tr().Body(
						app.Td().ColSpan(3).Text("Loading..."),
					)
				}),
                app.If(len(l.labResults)==0, func() app.UI {
					return app.Tr().Body(
						app.Td().ColSpan(3).Text("No Results Found."),
					)
				}),
				app.Range(l.labResults).Slice(func(i int) app.UI {
					row := l.labResults[i]
					reportDateFormatted := row.ReportDate.Format("02/01/2006")
					return app.Tr().Body(
						app.Td().Text(row.ID),
						app.Td().Text(reportDateFormatted),
						app.Td().Text(row.Subject),
					)
				}),
			),
		),
	)
}

func (l *LabResults) headerClass(field string) string {
    if l.SortBy != field {
        return "sortable"
    }
    if l.SortOrder == "asc" {
        return "sorted-asc"
    }
    return "sorted-desc"
}

func (l *LabResults) handleSort(ctx app.Context, field string) {
    if l.SortBy == field {
        // Toggle order if clicking the same column
        if l.SortOrder == "asc" { l.SortOrder = "desc" } else { l.SortOrder = "asc" }
    } else {
        l.SortBy = field
        l.SortOrder = "asc"
    }
    l.sortData()
    ctx.Update() // Trigger re-render
}

func (l *LabResults) sortData() {
    if l.SortBy == "" {
        return // No sort selected
    }

    sort.Slice(l.labResults, func(i, j int) bool {
        a := l.labResults[i]
        b := l.labResults[j]

        var less bool
        
        switch l.SortBy {
        case "ID":
            less = strings.ToLower(a.ID) < strings.ToLower(b.ID)
        case "Subject":
            less = strings.ToLower(a.Subject) < strings.ToLower(b.Subject)
        case "ReportDate":
            less = a.ReportDate.Before(b.ReportDate)
        default:
            return false
        }
        
        // Reverse order if descending
        if l.SortOrder == "desc" {
            return !less
        }
        return less
    })
}

func (l *LabResults) OnMount(ctx app.Context) {
    l.Loading = true
    ctx.Update()

    l.SortBy = "ReportDate"
    l.SortOrder = "desc"
    l.sortData()
    
    l.Loading = false
    ctx.Update()
}
