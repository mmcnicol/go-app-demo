package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
	"time"
)

type GopherBanner struct {
	app.Compo
	GopherDemographics *models.GopherDemographics
}

func NewGopherBanner(demographics *models.GopherDemographics) *GopherBanner {
	return &GopherBanner{
		GopherDemographics: demographics,
	}
}

func (gb *GopherBanner) Render() app.UI {

	if g.GopherDemographics == nil {
		return app.Div().Class("demographics-bar").Body(
			app.Span().Text("No Gopher Selected"),
		)
	}

	dobFormatted := gb.gopherDemographics.DateOfBirth.Format("02/01/2006")
	age := int(time.Since(gb.GopherDemographics.DateOfBirth).Hours() / 24 / 365)

	return app.Div().Class("demographics-bar").Body(
		app.Span().Text(gb.GopherDemographics.Name),
		app.Span().Text("DOB: "),
		app.Span().Text(dobFormatted),
		app.Span().Text(fmt.Sprintf(" (%d years old)", age)),
		//app.Span().Text(gb.GopherId),
	)
}
