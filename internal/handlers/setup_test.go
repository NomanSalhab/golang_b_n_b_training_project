package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func getRoutes() http.Handler {
	//* What to put in the Session
	gob.Register(models.Reservation{})

	// * Change This to true in production mode
	app.InProduction = false

	//? Creating A Logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour              //? 30Min for example for more secure auth
	session.Cookie.Persist = true                  //? if we want to logout when the session is closed we set Persist to False
	session.Cookie.SameSite = http.SameSiteLaxMode //? Cookie Structness
	session.Cookie.Secure = app.InProduction       //? in production we set it to true for https

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	// ? When in Development Mode UseCache is false
	app.UseCache = true

	repo := NewRepo(&app, nil)
	NewHandlers(repo)

	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer) // * For When The Application Panics

	//? We Commented NoSurf Because it requires a CSRFToken and we're not providing it in our tests
	// mux.Use(NoSurf)
	//?? NoSurf Prevents Any POST Request Without Proper CSRF Protection Token
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

	// return nil
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// * SetBaseCookie to make sure that the token it generates is available on a per page basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",              // * To Be Applied to the entire site
		Secure:   app.InProduction, // * because NOW we're not running on https
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and savesthe session on every request
func SessionLoad(next http.Handler) http.Handler {
	//? Provides Middleware which loads and saves session data for the current
	//? request and communicates the session token to and from the client in a cookie
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		// fmt.Println("index is:", i)
		name := filepath.Base(page)
		// fmt.Println("Page Name in CreateTemplateCache Function is:", name, "in Pages Variable:", pages)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
			myCache[name] = ts
			// ? Remember This Line Idiot
			// // return myCache, nil

		}

	}

	return myCache, nil
}
