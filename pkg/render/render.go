package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/nomansalhab/golang_training_project/pkg/config"
	"github.com/nomansalhab/golang_training_project/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplate sets the config for the new template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate Function
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// * Get The Template Cache From The AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// tc, err1 := CreateTemplateCache()
	// if err1 != nil {
	// 	fmt.Println("CreateTemplateCache Error:", err1)
	// 	log.Fatal(err1)
	// }
	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("could not get template from template cache")
		fmt.Println("!ok Error: could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err2 := t.Execute(buf, td)
	if err2 != nil {
		fmt.Println("Error at t.Execute", err2)
	}
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	// err := parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error parsing template:", err)
	// 	return
	// }

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
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
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
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
