package main

import (
	"encoding/json"
	"films-api/pkg/models"
	_ "films-api/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

/*
POST   /actor/<name>/<gender>/<birthday> :  добавление информации об актёре (имя, пол, дата рождения)
DELETE /actor/<actorid>    :  удаляет актёра по ID

GET    /task/<taskid>      :  возвращает одну задачу по её ID
GET    /task/              :  возвращает все задачи
DELETE /task/<taskid>      :  удаляет задачу по ID
GET    /tag/<tagname>      :  возвращает список задач с заданным тегом
GET    /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату

При добавлении фильма указываются его название
(не менее 1 и не более 150 символов),
описание (не более 1000 символов),
дата выпуска,
рейтинг (от 0 до 10) и список актёров:

POST /film/
*/

func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//key := r.Header.Get("API-Key")
		/*
				usr, err := app.users.Get(key)
			if err != nil {
				app.serverError(w, err)
				return
			}
			if key == usr.Api_key && usr.Role == "admin" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

		*/
		next.ServeHTTP(w, r)

		fmt.Println(r.URL)

	})
}

func (app *application) actor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.getActor(w, r)
		fmt.Println(r.URL)
	} else if r.Method == http.MethodPost {
		app.createActor(w, r)
	} else if r.Method == http.MethodDelete {
		app.deleteActor(w, r)
	} else if r.Method == http.MethodPut {
		app.changeActor(w, r)
	} else {
		http.Error(w, fmt.Sprintf("expect method GET, DELETE, POST or PUT at /actor, got %v", r.Method), http.StatusMethodNotAllowed)
	}

}

func (app *application) film(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.getFilm(w, r)
	} else if r.Method == http.MethodPost {
		app.createFilm(w, r)
	} else if r.Method == http.MethodDelete {
		app.deleteFilm(w, r)
	} else {
		http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /film, got %v", r.Method), http.StatusMethodNotAllowed)
	}

}

func (app *application) getActor(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	act, err := app.actors.Get_By_Id(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, act)

}

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

func (app *application) deleteActor(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	err = app.actors.Delete(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

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

func (app *application) getActors(w http.ResponseWriter, r *http.Request) {

	act, err := app.actors.Get_Actors()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderJSON(w, act)
}

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

func (app *application) deleteFilm(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))

	err = app.films.Delete(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

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
