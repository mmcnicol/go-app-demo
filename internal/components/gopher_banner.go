package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/state"
)

type GopherBanner struct {
	app.Compo
	GopherDemographics state.GopherDemographics
}

func NewGopherBanner(gopherDemographics state.GopherDemographics) *GopherBanner {
    return &GopherBanner{
        GopherDemographics: gopherDemographics,
    }
}

func (gb *GopherBanner) Render() app.UI {

	dobFormatted := gb.GopherDemographics.DateOfBirth.Format("02/01/2006")
	age := int(time.Since(gb.GopherDemographics.DateOfBirth).Hours() / 24 / 365)

	return app.Div().Class("gopher-banner").Body(
		app.Span().Text(gb.Name),
		app.Span().Text("DOB: "),
		app.Span().Text(dobFormatted),
		app.Span().Text(fmt.Sprintf(" (%d years old)", age)),
		//app.Span().Text(gb.GopherId),
	)
}
