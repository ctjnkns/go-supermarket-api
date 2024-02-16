package main

import "testing"

func Test_Init(t *testing.T) {
	var app Config
	err := app.init()
	if err != nil {
		t.Errorf("Failed init")
	}
}
