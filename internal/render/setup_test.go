package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//* What to put in the Session
	gob.Register(models.Reservation{})

	// * Change This to true in production mode
	testApp.InProduction = false

	//? Creating A Logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour              //? 30Min for example for more secure auth
	session.Cookie.Persist = true                  //? if we want to logout when the session is closed we set Persist to False
	session.Cookie.SameSite = http.SameSiteLaxMode //? Cookie Structness
	session.Cookie.Secure = false                  //? in production we set it to true for https

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// ** These Below Are All Fake Functions to Simulate a http.ResponseWriter
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}
