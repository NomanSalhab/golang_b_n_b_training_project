package main

import (
	"fmt"
	"testing"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// Nothing, All True
	default:
		t.Error(fmt.Sprintf("Type Isn't *chi.Mux, But Is: %T", v))
	}
}
