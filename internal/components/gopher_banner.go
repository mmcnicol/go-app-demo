package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
	"go-app-demo/internal/state"
	"time"
)

type GopherBanner struct {
	app.Compo
	GopherDemographics *models.GopherDemographics
	isMounted          bool
}

func NewGopherBanner() *GopherBanner {
	return &GopherBanner{}
}

func (gb *GopherBanner) OnMount(ctx app.Context) {
	app.Log("[GopherBanner OnMount] Initializing")
	
	ctx.ObserveState(state.GopherDemographicsKey, &gb.GopherDemographics).
		OnChange(func() {
			app.Log("[GopherBanner] Gopher demographics changed")
			gb.Refresh(ctx)
		})
	
	gb.isMounted = true
	gb.Refresh(ctx)
}

// Refresh updates the component
func (gb *GopherBanner) Refresh(ctx app.Context) {
	ctx.Update()
}

func (gb *GopherBanner) Render() app.UI {
	app.Log("[GopherBanner Render] Called, isMounted:", gb.isMounted, "demographics:", gb.GopherDemographics != nil)
	
	if !gb.isMounted {
		return app.Div().Class("demographics-bar loading")
	}

	if gb.GopherDemographics == nil {
		return app.Div().Class("demographics-bar empty").Body(
			app.Span().Class("no-selection").Text("No Gopher Selected"),
		)
	}

	dobFormatted := gb.GopherDemographics.DateOfBirth.Format("02/01/2006")
	age := int(time.Since(gb.GopherDemographics.DateOfBirth).Hours() / 24 / 365)

	return app.Div().Class("demographics-bar").Body(
		app.Div().Class("demographics-info").Body(
			app.Span().Class("gopher-name").Text(gb.GopherDemographics.Name),
			app.Span().Class("gopher-dob").Text(fmt.Sprintf("DOB: %s (%d years)", dobFormatted, age)),
			app.Span().Class("gopher-id").Text("ID: " + gb.GopherDemographics.GopherId),
		),
	)
}
