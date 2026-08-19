package web

import "errors"

// errNotImplemented marks a Phase 0 stub whose owning package has not
// implemented the real behavior yet.
//
// It was originally declared in generate.go. Now that Generate is implemented,
// the declaration lives here so that the remaining stub (presenter.go, owned by
// Package C) keeps compiling without that file needing to change.
//
// Delete this file once every Phase 0 stub is implemented.
var errNotImplemented = errors.New("not implemented")
