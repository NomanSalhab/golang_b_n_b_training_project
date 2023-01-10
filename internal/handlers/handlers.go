package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/driver"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/forms"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/helpers"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/render"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/repository"
	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/repository/dbrepo"
	"github.com/go-chi/chi"
)

// Repo the reposirory user by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a ner Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a ner Repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home Page Handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About Page Handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Reservation Page Handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cant't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation Handles The Posting of a Reservation Form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, _ /*ok*/ := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// if !ok {
	// 	m.App.Session.Put(r.Context(), "error", "can's get from session")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	startDate, errStart := http.ParseTime(sd)

	if errStart != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, errEnd := http.ParseTime(ed)

	if errEnd != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsEmail("email")
	form.MinLength("phone", 8)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Redirect(w, r, "/", http.StatusSeeOther)

		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err0 := m.DB.InsertReservation(reservation)
	if err0 != nil {
		m.App.Session.Put(r.Context(), "error", "cant't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err1 := m.DB.InsertRoomRestriction(restriction)
	if err1 != nil {
		m.App.Session.Put(r.Context(), "error", "cant't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals Page Handler Renders The Room Page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors Page Handler Renders The Room Page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

// Availability Page Handler Renders The Search Availability Page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability Page Handler Renders The Search Availability Page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{})

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-Jan-02" //"Mon, 01/02/06, 03:04PM" //
	startDate, startErr := time.Parse(layout, start)
	if startErr != nil {
		helpers.ServerError(w, startErr)
		return
	}
	endDate, endErr := time.Parse(layout, end)
	if endErr != nil {
		helpers.ServerError(w, endErr)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		//? No Availability
		// m.App.InfoLog.Println("No Availability!")
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

	// w.Write([]byte(fmt.Sprintf("Start Date Is: %s and End Date Is: %s", start, end)))
}

// * if this struct is only used in a method or a group of function => put it as close to it/them as possible
type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON Page Handler Handles The Search Availability Request & Sends A JSON Response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "      ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-Jan-02" //"Mon, 01/02/06, 03:04PM" //
	startDate, startError := time.Parse(layout, sd)
	if startError != nil {
		helpers.ServerError(w, startError)
		return
	}
	endDate, endError := time.Parse(layout, ed)
	if endError != nil {
		helpers.ServerError(w, endError)
		return
	}

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	// if roomIdError != nil {
	// 	helpers.ServerError(w, roomIdError)
	// 	return
	// }

	available, availableError := m.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	if availableError != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "      ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "     ")
	// if err != nil {
	// 	// log.Println(err)
	// 	helpers.ServerError(w, err)
	// 	return
	// }
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ChooseRoom displays a list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// m.App.Session.Get(r.Context(), "reservation")
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters and builds a sessionan variable and takes user ti /make-reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	layout := "2006-Jan-02" //"Mon, 01/02/06, 03:04PM" //
	startDate, startError := time.Parse(layout, sd)
	if startError != nil {
		helpers.ServerError(w, startError)
		return
	}
	endDate, endError := time.Parse(layout, ed)
	if endError != nil {
		helpers.ServerError(w, endError)
		return
	}
	room, _ := m.DB.GetRoomByID(res.RoomID)
	// if err != nil {
	// helpers.ServerError(w, err)
	// }

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

	log.Println(roomID, startDate, endDate)
}

// Availability Page Handler Renders The Search Availability Page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// Availability Page Handler Renders The Search Availability Page
// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// log.Println("can't Get Item From Session")
		m.App.ErrorLog.Println("can't Get Item From Session")
		m.App.Session.Put(r.Context(), "error", "Can't Get Reservation From Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
