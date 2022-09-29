package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

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
