package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// getHome renders the default page through the real router and template.
func getHome(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	app := newTestApp(t, parseProductionTemplate(t))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.Routes().ServeHTTP(response, request)
	return response
}

func TestHomeReturnsOKWithCompleteForm(t *testing.T) {
	response := getHome(t)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	body := response.Body.String()
	for _, fragment := range []string{
		`name="text"`,
		`name="banner"`,
		`value="standard"`,
		`value="shadow"`,
		`value="thinkertoy"`,
		`type="submit"`,
	} {
		if !strings.Contains(body, fragment) {
			t.Errorf("default page does not contain %q", fragment)
		}
	}
}

func TestHomeSelectsDefaultBanner(t *testing.T) {
	body := getHome(t).Body.String()

	standard := strings.Index(body, `value="standard"`)
	if standard < 0 {
		t.Fatal("standard banner option is missing")
	}
	// The checked attribute is rendered immediately after the value on the
	// selected radio, so a short window after it is enough to prove which
	// option the server marked as selected.
	window := body[standard:min(standard+80, len(body))]
	if !strings.Contains(window, "checked") {
		t.Error("standard banner is not selected by default")
	}

	if strings.Count(body, "checked") != 1 {
		t.Errorf("checked appears %d times, want exactly one selected banner",
			strings.Count(body, "checked"))
	}
}

func TestHomeStartsWithEmptyTextAndNoResult(t *testing.T) {
	body := getHome(t).Body.String()

	if !strings.Contains(body, `placeholder="Type something"></textarea>`) {
		t.Error("text area is not empty on the default page")
	}
	if strings.Contains(body, "columns") {
		t.Error("default page reports a rendered width without a result")
	}
}

func TestDefaultFormStateMatchesContract(t *testing.T) {
	form := defaultFormState()

	tests := []struct {
		field string
		got   string
		want  string
	}{
		{"Text", form.Text, ""},
		{"Banner", form.Banner, BannerStandard},
		{"Color", form.Color, "#ff0000"},
		{"Substring", form.Substring, ""},
		{"Alignment", form.Alignment, "left"},
		{"Width", form.Width, "80"},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("default %s = %q, want %q", test.field, test.got, test.want)
		}
	}
	if form.UseColor {
		t.Error("colouring is enabled by default, want it off")
	}
}

func TestClearReturnsTheDefaultPage(t *testing.T) {
	// Clear is defined as a fresh GET of "/", so it must produce exactly the
	// same response as a first visit rather than an HTML form reset.
	first := getHome(t).Body.String()
	second := getHome(t).Body.String()

	if first != second {
		t.Error("clearing the page does not reproduce the default page")
	}
}
