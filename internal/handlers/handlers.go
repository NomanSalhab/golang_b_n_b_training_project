package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/render"
)

// Repo the reposirory user by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a ner Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home Page Handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// sum := addValues(2, 2)
	// _, _ = fmt.Fprintf(w, fmt.Sprintf("This is The About Page and 2 + 2 is: %d", sum))

	//* We Will Grab The IP Address of The Visitor and Store it in the Seesion
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

// About Page Handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// sum := addValues(2, 2)
	// _, _ = fmt.Fprintf(w, fmt.Sprintf("This is The About Page and 2 + 2 is: %d", sum))

	// * We Will Perform Some Business Logic

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation Page Handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{})
}

// Generals Page Handler Renders The Room Page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors Page Handler Renders The Room Page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

// Availability Page Handler Renders The Search Availability Page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability Page Handler Renders The Search Availability Page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{})

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start Date Is: %s and End Date Is: %s", start, end)))
}

// * if this struct is only used in a method or a group of function => put it as close to it/them as possible
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON Page Handler Handles The Search Availability Request & Sends A JSON Response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type:", "application/json")
	w.Write(out)
}

// Availability Page Handler Renders The Search Availability Page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}
