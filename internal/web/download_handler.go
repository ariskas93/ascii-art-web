package web

import "net/http"

// handleDownload is a Phase 0 route placeholder. Package B replaces it with
// the validated, regenerated plain-text download workflow.
func (app *App) handleDownload(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
