package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/driver"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/handlers"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/helpers"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8016"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

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

func run() (*driver.DB, error) {

	//* What to put in the Session
	gob.Register(models.Reservation{})

	// * Change This to true in production mode
	app.InProduction = false

	//? Creating A Logger
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour              //? 30Min for example for more secure auth
	session.Cookie.Persist = true                  //? if we want to logout when the session is closed we set Persist to False
	session.Cookie.SameSite = http.SameSiteLaxMode //? Cookie Structness
	session.Cookie.Secure = app.InProduction       //? in production we set it to true for https

	app.Session = session

	//* Connect To Database

	log.Println("Connecting To Database!")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bed_n_breakfast user=postgres password=postgresqlgolangpass")
	if err != nil {
		log.Fatal("Couldn't Connected To Database!")
	}
	log.Println("Connected To Database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	// ? When in Development Mode UseCache is false
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
