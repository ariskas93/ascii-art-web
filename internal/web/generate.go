package web

import "errors"

var errNotImplemented = errors.New("not implemented")

// Generate runs the shared banner lookup, normalization, rendering, and
// alignment workflow. Package A implements the workflow after Phase 0.
func (app *App) Generate(in GenerationInput) (GeneratedASCII, error) {
	return GeneratedASCII{}, errNotImplemented
}
