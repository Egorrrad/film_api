package main

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)
import _ "films-api/docs"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/actor", app.actor)
	mux.HandleFunc("/api/actors", app.getActors)

	mux.HandleFunc("/api/film", app.film)
	mux.HandleFunc("/api/films/sort", app.getFilms)
	mux.HandleFunc("/api/film/search", app.searchFilm)

	mux1 := http.NewServeMux()
	mux1.Handle("/api/", app.apiKeyMiddleware(mux))

	config := httpSwagger.URL("http://localhost:4000/swagger/doc.json")
	mux1.HandleFunc("/swagger/", httpSwagger.Handler(config))

	return mux1
}
