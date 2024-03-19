package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// Access is denied
func (app *application) checkAdmin(w http.ResponseWriter, api_key string) bool {
	result, _ := app.users.IsAdmin(api_key)

	if !result {
		http.Error(w, "Access is denied", http.StatusForbidden)
		return false
	}
	return true
}

func (app *application) printInfoRequest(r *http.Request) {
	app.infoLog.Printf("%s %s", r.Method, r.URL)
}
