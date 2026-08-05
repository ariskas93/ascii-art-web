package web

import "net/http"

// Routes returns the application's complete HTTP handler.
func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.handleHome)
	mux.HandleFunc("POST /ascii-art", app.handleGenerate)
	mux.HandleFunc("POST /ascii-art/download", app.handleDownload)
	return mux
}
