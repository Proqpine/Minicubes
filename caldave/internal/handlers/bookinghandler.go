package handlers

import (
	"net/http"
)

func BookingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "booking.html", "Hello World")
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	})
}
