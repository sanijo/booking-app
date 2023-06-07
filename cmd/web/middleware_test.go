package main

import (
	"fmt"
	"net/http"
	"testing"

)

func TestNoSurf(t *testing.T) {
    // Create a dummy next handler for testing
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	h := NoSurf(nextHandler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
    // Create a dummy next handler for testing
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	h := SessionLoad(nextHandler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
