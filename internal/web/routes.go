package web

import "net/http"

// Routes returns the application's complete HTTP handler.
func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.handleHome)
	mux.HandleFunc("POST /ascii-art", app.handleGenerate)
	mux.HandleFunc("POST /ascii-art/download", app.handleDownload)
	mux.HandleFunc("GET /static/style.css", app.handleStyleCSS)
	mux.HandleFunc("GET /static/app.js", app.handleAppJS)
	return mux
}

func (app *App) handleStyleCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/style.css")
}

func (app *App) handleAppJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/app.js")
}
