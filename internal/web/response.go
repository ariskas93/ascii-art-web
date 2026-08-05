package web

import (
	"bytes"
	"log/slog"
	"net/http"
)

// renderPage buffers template output so execution failures can become a clean
// error response instead of a partially written successful page.
func (app *App) renderPage(w http.ResponseWriter, status int, data PageData) {
	var body bytes.Buffer
	if err := app.pageTemplate.ExecuteTemplate(&body, "page", data); err != nil {
		app.logger.Error("execute page template", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(body.Bytes()); err != nil {
		app.logger.Error("write page response", slog.Any("error", err))
	}
}
