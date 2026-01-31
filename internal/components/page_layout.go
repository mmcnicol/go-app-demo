package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/models"
)

type PageLayout struct {
    app.Compo
    ApplicationBanner *ApplicationBanner
    LeftNavigation    *LeftNavigation
    GopherBanner      *GopherBanner
    Body              app.UI // This acts like a "slot" or "children"
    PageFooter        *PageFooter
}

func NewPageLayout(body app.UI) *PageLayout {
    return &PageLayout{
        ApplicationBanner: NewApplicationBanner("Application Banner", nil, nil),
        LeftNavigation:    NewLeftNavigation([]models.NavItem{}),
        GopherBanner:      NewGopherBanner(&models.GopherDemographics{}),
        Body:              body,
        PageFooter:        NewPageFooter("Page Footer"),
    }
}

// Alternative constructor that accepts all components
func NewPageLayoutWithComponents(appBanner *ApplicationBanner, nav *LeftNavigation, 
    gopherBanner *GopherBanner, body app.UI, footer *PageFooter) *PageLayout {
    return &PageLayout{
        ApplicationBanner: appBanner,
        LeftNavigation:    nav,
        GopherBanner:      gopherBanner,
        Body:              body,
        PageFooter:        footer,
    }
}

func (p *PageLayout) Render() app.UI {
    return app.Div().Class("app-container").Body(
        // Use pointer to applicationBanner (which implements app.UI)
        p.ApplicationBanner,
        app.Div().Class("app-body").Body(
            // Include left navigation if it exists
            app.If(p.LeftNavigation != nil,
                func() app.UI {
                    return p.LeftNavigation
                },
            ),
            app.Div().Class("main-content").Body(
                // Use pointer to gopherBanner
                p.GopherBanner,
                app.Div().Class("content-area").Body(
                    p.Body,
                ),
                // Use pointer to pageFooter
                p.PageFooter,
            ),
        ),
    )
}

// Setter methods to update components
func (p *PageLayout) SetBody(ctx app.Context, body app.UI) {
    p.Body = body
    ctx.Update()
}

func (p *PageLayout) SetLeftNavigation(ctx app.Context, nav *LeftNavigation) {
    p.LeftNavigation = nav
    ctx.Update()
}

func (p *PageLayout) SetApplicationBanner(ctx app.Context, banner *ApplicationBanner) {
    p.ApplicationBanner = banner
    ctx.Update()
}

func (p *PageLayout) SetGopherBanner(ctx app.Context, banner *GopherBanner) {
    p.GopherBanner = banner
    ctx.Update()
}

func (p *PageLayout) SetPageFooter(ctx app.Context, footer *PageFooter) {
    p.PageFooter = footer
    ctx.Update()
}
