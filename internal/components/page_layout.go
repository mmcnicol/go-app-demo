package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/components"
	//"go-app-demo/internal/models"
)

type PageLayout struct {
    app.Compo
    applicationBanner ApplicationBanner
    leftNavigation LeftNavigation
    gopherBanner GopherBanner
    Body app.UI // This acts like a "slot" or "children"
    pageFooter PageFooter
}

func (p *PageLayout) Render() app.UI {
    return app.Div().Class("app-container").Body(
        p.applicationBanner,
        app.Div().Class("app-body").Body(
            app.Div().Class("main-content").Body(
                p.gopherBanner,
                app.Div().Class("content-area").Body(
                    p.Body,
                ),
                p.pageFooter,
            ),
        ),
    )
}
