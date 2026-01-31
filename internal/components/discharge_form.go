package components

import (
	//"encoding/json"
	//"fmt"
	//"net/http"
	//"strings"
	//"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go-app-demo/internal/models"
)

type DischargeForm struct {
	app.Compo
	patientID string
	dischargeForm models.DischargeForm
}

// OnMount is called when the component is first loaded
func (d *DischargeForm) OnMount(ctx app.Context) {
	d.dischargeForm = models.DischargeForm{
		ID:   "MRN-8829",
		Name: "John Doe",
	}
}

/*
func (d *DischargeForm) OnNav(ctx app.Context) {
    // 1. Grab the ID from the URL
    d.patientID = ctx.Page().URL().Path[11:] // Simple string slicing or use a regex
    
    // 2. Simulate fetching the full record based on that ID
    // In a real app, you might call an API here
    d.loadPatientData(d.patientID)
}
*/

func (d *DischargeForm) OnNav(ctx app.Context) {
    
	/*
	// Retrieve the object from state
    ctx.GetState("current-patient", &d.patient)
    ctx.Update()
	*/
}

/*
func (d *DischargeForm) loadPatientData(id string) {
    
	// Simulation: would normally be a database lookup
    if id == "MRN-8829" {
        d.dischargeForm = models.DischargeForm{
			ID:   "MRN-8829",
			Name: "John Doe",
		}
    }
    ctx.Update()
}
*/

func (d *DischargeForm) Render() app.UI {
	return app.Div().Class("container").Body(
		app.H1().Text("Clinical Portal"),
		app.P().Text("Patient Discharge System"),
		
		app.Div().Class("form-group").Body(
			app.Label().Text("Patient Name: "),
			app.Input().
				Type("text").
				Value(d.dischargeForm.Name),

			app.Label().Text("Diagnosis: "),
			app.Input().
				Type("text").
				Value(d.dischargeForm.Diagnosis).
				OnChange(d.ValueTo(&d.dischargeForm.Diagnosis)), // Updates data automatically!

		),

		app.Div().Class("summary").Body(
			app.H3().Text("Discharge Summary"),
			app.P().Body(
				app.B().Text("Ready to discharge: "),
				app.Text(d.dischargeForm.Name),
				app.Text(" (ID: "+d.dischargeForm.ID+")"),
				app.Text(" (Diagnosis: "+d.dischargeForm.Diagnosis+")"),
			),
		),
	)
}
