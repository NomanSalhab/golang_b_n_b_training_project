package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name   string
	url    string
	method string
	// params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET" /*[]postData{},*/, http.StatusOK},
	{"about", "/about", "GET" /*[]postData{},*/, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET" /*[]postData{},*/, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET" /*[]postData{},*/, http.StatusOK},
	{"search-availability", "/search-availability", "GET" /*[]postData{},*/, http.StatusOK},
	{"contact", "/contact", "GET" /*[]postData{},*/, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", /*[]postData{},*/ http.StatusOK},

	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2022-10-11"},
	// 	{key: "end", value: "2022-10-13"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2022-10-11"},
	// 	{key: "end", value: "2022-10-13"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Noman"},
	// 	{key: "last_name", value: "Salhab"},
	// 	{key: "email", value: "nomansalhab@gmail.com"},
	// 	{key: "phone", value: "0992008569"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//* Creating a Test Server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() //? To Close The Test Server When The Test Function Is Finished

	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("For %s, Expedted %d But Got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation Handler returned wrong response code: got: %d , wanted: %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation Handler returned wrong response code: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation Handler returned wrong response code: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "Noman")
	postedData.Add("last_name", "Salhab")
	postedData.Add("email", "noman@gmail.com")
	postedData.Add("phone", "0992008516")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		// t.Errorf("PostReservation Handler returned wrong response code: got: %d , wanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for midding post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler returned wrong response code for missing post body: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date

	// reqBody = "start_date=invalid"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData.Set("start_date", "invalid")
	postedData.Set("end_date", "2050-01-02")
	postedData.Set("first_name", "Noman")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler returned wrong response code for parseng start date: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date

	// reqBody = "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData.Set("start_date", "2050-01-01")
	postedData.Set("end_date", "invalid")
	postedData.Set("first_name", "Noman")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler returned wrong response code for parseng end date: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid room id

	// reqBody = "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-05")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	postedData.Set("start_date", "2050-01-01")
	postedData.Set("end_date", "2050-01-02")
	postedData.Set("first_name", "Noman")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler returned wrong response code for invalid room id: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid invalid data

	// reqBody = "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-05")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=N")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData.Set("start_date", "2050-01-01")
	postedData.Set("end_date", "2050-01-02")
	postedData.Set("first_name", "N")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler returned wrong response code for invalid data: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid faliure inserting reservation into database

	// reqBody = "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-05")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	postedData.Set("start_date", "2050-01-01")
	postedData.Set("end_date", "2050-01-02")
	postedData.Set("first_name", "Noman")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler failed when trying to fail inserting reservation: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid faliure inserting restriction into database

	// reqBody = "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-05")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Noman")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Salhab")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=noman@gmail.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0992008569")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	postedData.Set("start_date", "2050-01-01")
	postedData.Set("end_date", "2050-01-02")
	postedData.Set("first_name", "Noman")
	postedData.Set("last_name", "Salhab")
	postedData.Set("email", "noman@gmail.com")
	postedData.Set("phone", "0992008516")
	postedData.Set("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//* Tells The Web Server that The Request is of Type (Form Post)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation Handler failed when trying to fail inserting restriction: got: %d , wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepositoryAvailabilityJson(t *testing.T) {
	// first case: rooms are not available
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// Get Context With Session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//set The Request Header
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	// Make Handler HandlerFunc

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// Get Response Reqorder
	rr := httptest.NewRecorder()

	// Make The Request To Our Handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
