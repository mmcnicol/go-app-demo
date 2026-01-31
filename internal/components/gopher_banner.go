package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
	"time"
)

type GopherBanner struct {
	app.Compo
	GopherDemographics models.GopherDemographics
}

func NewGopherBanner(gopherDemographics models.GopherDemographics) *GopherBanner {
    return &GopherBanner{
        GopherDemographics: gopherDemographics,
    }
}

func (gb *GopherBanner) Render() app.UI {

	dobFormatted := gb.GopherDemographics.DateOfBirth.Format("02/01/2006")
	age := int(time.Since(gb.GopherDemographics.DateOfBirth).Hours() / 24 / 365)

	return app.Div().Class("demographics-bar").Body(
		app.Span().Text(gb.GopherDemographics.Name),
		app.Span().Text("DOB: "),
		app.Span().Text(dobFormatted),
		app.Span().Text(fmt.Sprintf(" (%d years old)", age)),
		//app.Span().Text(gb.GopherId),
	)
}
