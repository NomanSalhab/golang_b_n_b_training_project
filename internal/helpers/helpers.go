package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
)

var app *config.AppConfig

// NewHelpers Sets Up AppConfig For Helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error With A Status Of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
