package web

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAppValidatesDependencies(t *testing.T) {
	pageTemplate := parseProductionTemplate(t)
	registry := testBannerRegistry()
	logger := testLogger()

	tests := []struct {
		name      string
		template  *template.Template
		registry  BannerRegistry
		logger    *slog.Logger
		wantError string
	}{
		{
			name:      "missing template",
			registry:  registry,
			logger:    logger,
			wantError: "page template is required",
		},
		{
			name:      "missing definition",
			template:  template.Must(template.New("incomplete").Parse(`{{define "page"}}page{{end}}`)),
			registry:  registry,
			logger:    logger,
			wantError: "missing template definition",
		},
		{
			name:      "missing logger",
			template:  pageTemplate,
			registry:  registry,
			wantError: "logger is required",
		},
		{
			name:      "missing standard banner",
			template:  pageTemplate,
			registry:  registryWithout(BannerStandard),
			logger:    logger,
			wantError: `banner "standard" is required`,
		},
		{
			name:      "missing shadow banner",
			template:  pageTemplate,
			registry:  registryWithout(BannerShadow),
			logger:    logger,
			wantError: `banner "shadow" is required`,
		},
		{
			name:      "missing thinkertoy banner",
			template:  pageTemplate,
			registry:  registryWithout(BannerThinkertoy),
			logger:    logger,
			wantError: `banner "thinkertoy" is required`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewApp(test.template, test.registry, test.logger)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("NewApp() error = %v, want error containing %q", err, test.wantError)
			}
		})
	}
}

func TestNewAppCopiesBannerRegistry(t *testing.T) {
	registry := testBannerRegistry()
	app, err := NewApp(parseProductionTemplate(t), registry, testLogger())
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}

	registry[BannerStandard] = map[rune][]string{'X': {"changed"}}
	delete(registry, BannerShadow)

	if _, ok := app.banners[BannerStandard][' ']; !ok {
		t.Fatal("replacing a caller registry entry changed the app registry")
	}
	if _, ok := app.banners[BannerShadow]; !ok {
		t.Fatal("deleting a caller registry entry changed the app registry")
	}
}

func TestRoutesPhaseZeroContract(t *testing.T) {
	app := newTestApp(t, parseProductionTemplate(t))
	handler := app.Routes()

	tests := []struct {
		name        string
		method      string
		path        string
		wantStatus  int
		wantAllowed string
	}{
		{name: "home", method: http.MethodGet, path: "/", wantStatus: http.StatusOK},
		{name: "unknown", method: http.MethodGet, path: "/unknown", wantStatus: http.StatusNotFound},
		{name: "wrong home method", method: http.MethodPost, path: "/", wantStatus: http.StatusMethodNotAllowed, wantAllowed: http.MethodGet},
		{name: "wrong generate method", method: http.MethodGet, path: "/ascii-art", wantStatus: http.StatusMethodNotAllowed, wantAllowed: http.MethodPost},
		{name: "wrong download method", method: http.MethodGet, path: "/ascii-art/download", wantStatus: http.StatusMethodNotAllowed, wantAllowed: http.MethodPost},
		{name: "generate placeholder", method: http.MethodPost, path: "/ascii-art", wantStatus: http.StatusNotImplemented},
		{name: "download placeholder", method: http.MethodPost, path: "/ascii-art/download", wantStatus: http.StatusNotImplemented},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.wantAllowed != "" && !headerContainsToken(response.Header().Get("Allow"), test.wantAllowed) {
				t.Errorf("Allow = %q, want it to contain %q", response.Header().Get("Allow"), test.wantAllowed)
			}
		})
	}
}

func TestHomeReturnsCompleteHTMLShell(t *testing.T) {
	app := newTestApp(t, parseProductionTemplate(t))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app.Routes().ServeHTTP(response, request)

	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
	}
	for _, fragment := range []string{
		"<!doctype html>",
		`action="/ascii-art"`,
		`href="/static/style.css"`,
		`src="/static/app.js"`,
	} {
		if !strings.Contains(response.Body.String(), fragment) {
			t.Errorf("response body does not contain %q", fragment)
		}
	}
}

func TestTemplateFailureReturnsCleanInternalServerError(t *testing.T) {
	failingTemplate := template.Must(template.New("index.html").Parse(failingPageTemplate))
	app := newTestApp(t, failingTemplate)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if strings.Contains(response.Body.String(), "partial page") {
		t.Fatalf("response leaked partial template output: %q", response.Body.String())
	}
	if response.Body.String() != http.StatusText(http.StatusInternalServerError)+"\n" {
		t.Errorf("body = %q, want a generic internal error", response.Body.String())
	}
}

func parseProductionTemplate(t *testing.T) *template.Template {
	t.Helper()
	pageTemplate, err := template.ParseFiles("../../templates/index.html")
	if err != nil {
		t.Fatalf("parse production template: %v", err)
	}
	return pageTemplate
}

func newTestApp(t *testing.T, pageTemplate *template.Template) *App {
	t.Helper()
	app, err := NewApp(pageTemplate, testBannerRegistry(), testLogger())
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	return app
}

func testBannerRegistry() BannerRegistry {
	return BannerRegistry{
		BannerStandard:   {' ': {"standard"}},
		BannerShadow:     {' ': {"shadow"}},
		BannerThinkertoy: {' ': {"thinkertoy"}},
	}
}

func registryWithout(name string) BannerRegistry {
	registry := testBannerRegistry()
	delete(registry, name)
	return registry
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func headerContainsToken(header string, token string) bool {
	for _, value := range strings.Split(header, ",") {
		if strings.TrimSpace(value) == token {
			return true
		}
	}
	return false
}

const failingPageTemplate = `
{{define "page"}}partial page {{.Missing}}{{end}}
{{define "input_controls"}}{{end}}
{{define "generate_action"}}{{end}}
{{define "output_box"}}{{end}}
{{define "align_controls"}}{{end}}
{{define "form_error"}}{{end}}
{{define "clear_action"}}{{end}}
{{define "download_action"}}{{end}}
{{define "color_controls"}}{{end}}
{{define "width_field"}}{{end}}
{{define "theme_control"}}{{end}}
`
