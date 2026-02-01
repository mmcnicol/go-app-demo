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
	isMounted     bool // Managed internally
	ctx app.Context
}

/*
func NewGopherBanner(demographics *models.GopherDemographics) *GopherBanner {
	return &GopherBanner{
		GopherDemographics: demographics,
	}
}
*/

func NewGopherBanner() *GopherBanner {
	return &GopherBanner{}
}

func (gb *GopherBanner) OnMount(ctx app.Context) {
	gb.ctx = ctx

	// Subscribe to all relevant state
	ctx.ObserveState(state.GopherDemographicsObservable).Value(&gb.GopherDemographics)

	// Update component when state changes
    ctx.ObserveState(state.GopherDemographicsObservable).OnChange(ctx.Update)

	// Read state
	ctx.ObserveState(state.GopherDemographicsObservable).Value(&gb.GopherDemographics)

	gb.isMounted = true
}

func (gb *GopherBanner) Render() app.UI {
	app.Log("[GopherBanner Render] Called")
	app.Log("[GopherBanner Render] Called, isMounted:", n.isMounted)
	
	// Don't render anything until mounted
	if !n.isMounted {
		app.Log("[GopherBanner Render] Not mounted yet, returning loading/empty")
		return app.Div().Class("demographics-bar loading")
	}

	//var demographics *models.GopherDemographics
    //ctx.ObserveState(state.GopherDemographicsObservable).Value(&demographics)

	if gb.GopherDemographics == nil {
		return app.Div().Class("demographics-bar").Body(
			app.Span().Text("No Gopher Selected"),
		)
	}

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
