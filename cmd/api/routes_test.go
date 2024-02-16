package main

import (
	"net/http"
	"testing"
)

func Test_Routes(t *testing.T) {
	var app Config
	h := app.routes()

	switch v := h.(type) {
	case http.Handler:
		//do nothing, passes
	default:
		t.Errorf("received type is not http.Handler, %T", v)
	}
}
