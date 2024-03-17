package main

import (
	"net/http"
)

// Middleware для проверки API ключа
func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка авторизации по API ключу
		if r.Header.Get("API_KEY") != "your_api_key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Прохождение авторизации
		next.ServeHTTP(w, r)
	})
}

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/actor", app.actor)
	mux.HandleFunc("/actors", app.getActors)

	mux.HandleFunc("/film", app.film)
	mux.HandleFunc("/films/sort", app.getFilms)
	mux.HandleFunc("/film/search", app.searchFilm)

	mux1 := http.NewServeMux()
	mux1.Handle("/api/", authHandler(mux))

	return mux1
}
