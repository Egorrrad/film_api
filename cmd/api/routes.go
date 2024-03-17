package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux1 := http.NewServeMux()

	mux := http.NewServeMux()
	mux.HandleFunc("actor", app.actor)
	mux.HandleFunc("/actors", app.getActors)

	mux.HandleFunc("/film", app.film)
	mux.HandleFunc("/films/sort", app.getFilms)
	mux.HandleFunc("/film/search", app.searchFilm)

	mux1.Handle("/api/", apiKeyMiddleware(mux))

	return mux1
}
