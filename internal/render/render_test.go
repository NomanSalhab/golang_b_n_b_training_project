package render

import (
	"net/http"
	"testing"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123456789")

	result := AddDefaultData(&td, r)

	if result.Flash != "123456789" {
		t.Error("Flash Value of 123456789 Is'nt Found In Session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err1 := getSession()
	if err1 != nil {
		t.Error(err1)
	}

	var ww myWriter

	err = Template(&ww, r, "home.page.html", &models.TemplateData{})

	if err != nil {
		t.Error("Error Writing Template To Browser!")
	}

	err = Template(&ww, r, "random.page.html", &models.TemplateData{})

	if err == nil {
		t.Error("Rendered Template That Doesn't Exist")
	}
}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/random-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) //! X-Session is case sensetive
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}
