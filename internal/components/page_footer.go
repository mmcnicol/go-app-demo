package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// Usage:
// &PageFooter{}
type PageFooter struct {
	app.Compo
	CopyrightText string
}

func NewPageFooter(copyrightText string) *PageFooter {
    return &PageFooter{
        CopyrightText: copyrightText,
    }
}

func (p *PageFooter) Render() app.UI {
    return app.Footer().Class("page-footer").Body(
        app.Div().Class("footer-content").Body(
            app.Span().Text(p.CopyrightText),
        ),
    )
}
