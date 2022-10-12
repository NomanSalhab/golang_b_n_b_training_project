package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/handlers"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8016"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting Applicationon Port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err1)
	}
}

func run() error {

	//* What to put in the Session
	gob.Register(models.Reservation{})

	// * Change This to true in production mode
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour              //? 30Min for example for more secure auth
	session.Cookie.Persist = true                  //? if we want to logout when the session is closed we set Persist to False
	session.Cookie.SameSite = http.SameSiteLaxMode //? Cookie Structness
	session.Cookie.Secure = app.InProduction       //? in production we set it to true for https

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	// ? When in Development Mode UseCache is false
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	return nil
}
