package main

import (
	"encoding/json"
	"films-api/pkg/models"
	_ "films-api/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

// Middleware для проверки API ключа

func (app *application) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("API-Key")

		usr, err := app.users.Get(key)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		if key == usr.Api_key {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

	})
}

func (app *application) actor(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("API-Key")
	if r.Method == http.MethodGet {
		app.getActor(w, r)
	} else {
		if !app.checkAdmin(w, key) {
			return
		}

		if r.Method == http.MethodPost {
			app.createActor(w, r)
		} else if r.Method == http.MethodDelete {
			app.deleteActor(w, r)
		} else if r.Method == http.MethodPut {
			app.changeActor(w, r)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE, POST or PUT at /actor, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}

}

func (app *application) film(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("API-Key")
	if r.Method == http.MethodGet {
		app.getFilm(w, r)
	} else {
		if !app.checkAdmin(w, key) {
			return
		}
		if r.Method == http.MethodPost {
			app.createFilm(w, r)
		} else if r.Method == http.MethodDelete {
			app.deleteFilm(w, r)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /film, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}

}

//	@Summary		Получение информации об актере
//	@Description	Получение информации о конкретном актере по id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Actor ID"
//	@Success		200	{object}	models.Actor
//
// @Failure      400  {string}  string    "error"
// @Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/actor [get]
func (app *application) getActor(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	if err != nil {
		app.serverError(w, err)
		return
	}
	act, err := app.actors.Get_By_Id(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, act)

}

// createActor godoc
//
//	@Summary		Добавление информации об актере
//	@Description	Добавление информации об актёре (имя, пол, дата рождения) через JSON
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			account	body		models.Actor	true	"Add actor"
//	@Success		200		{integer}	int				id
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500		{string}	string			"server error"
//	@Router			/api/actor [post]
func (app *application) createActor(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var act models.Actor
	err := decoder.Decode(&act)
	if err != nil {
		app.serverError(w, err)
	}

	id, err := app.actors.Insert(act.Name, act.Gender, act.Birthday)
	if err != nil {
		app.serverError(w, err)
		return
	}

	type answer struct {
		Id int `json:"id"`
	}
	app.renderJSON(w, answer{Id: id})

}

// deleteActor godoc
//
//	@Summary		Удаление информации об актере
//	@Description	Удаление информации об актёре по id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Actor ID"
//	@Success		200	{string}	string	"ok"
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/actor [delete]
func (app *application) deleteActor(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	err = app.actors.Delete(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

// changeActor godoc
//
//	@Summary		Изменение информации об актере
//	@Description	Возможно изменить любую информацию об актёре по его id, как частично, так и полностью
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			account	body		models.Actor	true	"Change actor"
//	@Success		200		{integer}	int				id
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500		{string}	string			"server error"
//	@Router			/api/actor [put]
func (app *application) changeActor(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var act models.Actor
	err := decoder.Decode(&act)
	if err != nil {
		app.serverError(w, err)
	}

	err = app.actors.Change(act.Name, act.Gender, act.Birthday, act.Id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	type answer struct {
		Id int `json:"id"`
	}
	app.renderJSON(w, answer{Id: act.Id})

}

// getActors godoc
//
//	@Summary		Получение списка актёров
//	@Description	Получение списка актёров, для каждого актёра выдаётся также список фильмов с его участием
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//
// Param			actor	body		models.Actor	true	"Change actor"
//
//	@Success		200	{object}	models.Actors
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/actors [get]
func (app *application) getActors(w http.ResponseWriter, r *http.Request) {

	act, err := app.actors.Get_Actors()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, act)
}

// getFilm godoc
//
//	@Summary		Получение фильма по id
//	@Description	Получение фильма по его id
//	@Tags			film
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Film ID"
//	@Success		200	{object}	models.Film
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/film [get]
func (app *application) getFilm(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	fil, err := app.films.Get_By_Id(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, fil)

}

// createFilm godoc
//
//		@Summary		Добавление фильма
//		@Description	При добавлении фильма указываются его название (не менее 1 и не более 150 символов), описание (не более 1000 символов), дата выпуска, рейтинг (от 0 до 10) и список актёров:
//		@Tags			film
//		@Accept			json
//		@Produce		json
//	 @Param			film	body		models.Film	true	"Add film"
//		@Success		200	{integer}	int		id
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/film [post]
func (app *application) createFilm(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var fil models.Film
	err := decoder.Decode(&fil)
	if err != nil {
		app.serverError(w, err)
	}

	id, err := app.films.Insert(fil.Name, fil.Description, fil.Date, fil.Rating, fil.Actors)
	if err != nil {
		app.serverError(w, err)
		return
	}

	type answer struct {
		Id int `json:"id"`
	}
	app.renderJSON(w, answer{Id: id})

}

// not works
func (app *application) changeFilm(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	app.infoLog.Println(id)
	app.infoLog.Println(params)

	actor1 := "actor"
	js, err := json.Marshal(actor1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

// deleteFilm godoc
//
//	@Summary		Удаление фильма
//	@Description	Удаление фильма по его id
//	@Tags			film
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Film ID"
//	@Success		200	{integer}	int		id
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/film [delete]
func (app *application) deleteFilm(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	err = app.films.Delete(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

// getFilms godoc
//
//	@Summary		Получение фильма по id
//	@Description	Получение фильма по его id
//	@Tags			film
//	@Accept			json
//	@Produce		json
//	@Param			by	path		string	true	"Sorted By"
//	@Success		200	{object}	models.Films
//
// Failure      400  {string}  string    "error"
// Failure      404  {string}  string    "error"
//
//	@Failure		500	{string}	string	"server error"
//	@Router			/api/films [get]
func (app *application) getFilms(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	orderBy := params.Get("by")

	films, err := app.films.Get_Films(orderBy)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, films)
}

// not works
func (app *application) searchFilm(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	app.infoLog.Println(id)
	app.infoLog.Println(params)

	actor1 := "actor"
	js, err := json.Marshal(actor1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
