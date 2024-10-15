package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"time"
)

//go:embed *.html
var templateFiles embed.FS

var tpl *template.Template

type Calen struct {
	Day   time.Weekday
	Month time.Month
	Year  int
}

func init() {
	var err error
	tpl, err = template.ParseFS(templateFiles, "*.html")
	if err != nil {
		panic(err)
	}
}

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "index.html", "Hello World")
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	})
}
