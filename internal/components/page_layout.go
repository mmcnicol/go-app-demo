package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
    "go-app-demo/internal/models"
)

type PageLayout struct {
    app.Compo
    applicationBanner *ApplicationBanner
    leftNavigation    *LeftNavigation
    gopherBanner      *GopherBanner
    Body              app.UI // This acts like a "slot" or "children"
    pageFooter        *PageFooter
}

func NewPageLayout(body app.UI) *PageLayout {
    return &PageLayout{
        applicationBanner: NewApplicationBanner("Application Banner"),
        leftNavigation:    NewLeftNavigation([]models.NavItem{}, "", ""),
        gopherBanner:      NewGopherBanner(models.GopherDemographics{}),
        Body:              body,
        pageFooter:        NewPageFooter("Page Footer"),
    }
}

// Alternative constructor that accepts all components
func NewPageLayoutWithComponents(appBanner *ApplicationBanner, nav *LeftNavigation, 
    gopherBanner *GopherBanner, body app.UI, footer *PageFooter) *PageLayout {
    return &PageLayout{
        applicationBanner: appBanner,
        leftNavigation:    nav,
        gopherBanner:      gopherBanner,
        Body:              body,
        pageFooter:        footer,
    }
}

func (p *PageLayout) Render() app.UI {
    return app.Div().Class("app-container").Body(
        // Use pointer to applicationBanner (which implements app.UI)
        p.applicationBanner,
        app.Div().Class("app-body").Body(
            // Include left navigation if it exists
            app.If(p.leftNavigation != nil,
                func() app.UI {
                    return p.leftNavigation
                },
            ),
            app.Div().Class("main-content").Body(
                // Use pointer to gopherBanner
                p.gopherBanner,
                app.Div().Class("content-area").Body(
                    p.Body,
                ),
                // Use pointer to pageFooter
                p.pageFooter,
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
    p.leftNavigation = nav
    ctx.Update()
}

func (p *PageLayout) SetApplicationBanner(ctx app.Context, banner *ApplicationBanner) {
    p.applicationBanner = banner
    ctx.Update()
}

func (p *PageLayout) SetGopherBanner(ctx app.Context, banner *GopherBanner) {
    p.gopherBanner = banner
    ctx.Update()
}

func (p *PageLayout) SetPageFooter(ctx app.Context, footer *PageFooter) {
    p.pageFooter = footer
    ctx.Update()
}
