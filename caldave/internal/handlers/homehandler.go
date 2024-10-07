package handlers

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed *.html
var templateFiles embed.FS

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.ParseFS(templateFiles, "booking.html")
	if err != nil {
		panic(err)
	}
}

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Execute(w, "Hello World")
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	})
}
