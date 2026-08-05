package web

import "net/http"

func (app *App) handleHome(w http.ResponseWriter, _ *http.Request) {
	app.renderPage(w, http.StatusOK, PageData{Form: defaultFormState()})
}

// handleGenerate is a Phase 0 route placeholder. Package A replaces it with
// the parse, validate, generate, present, and render workflow.
func (app *App) handleGenerate(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func defaultFormState() FormState {
	return FormState{
		Banner:    BannerStandard,
		Color:     "#ff0000",
		Alignment: "left",
		Width:     "80",
	}
}
