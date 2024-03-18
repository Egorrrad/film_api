package main

import (
	"database/sql"
	"films-api/pkg/models/postgresql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	actors   *postgresql.ActorModel
	films    *postgresql.FilmModel
	users    *postgresql.UserModel
}

//	@title			Film App API
//	@version		1.0
//	@description	Это API для сервера фильмотеки

//	@host		localhost:4000
//	@BasePath	/api

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
// description Type 'Bearer TOKEN' to correctly set the API Key

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	//dsn := flag.String("dsn", "api:web00top@/snippetbox?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	connStr := "user=api_tester password=testing dbname=film_api sslmode=disable"
	db, err := openDB(connStr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// И добавляем его в зависимостях нашего
	// веб-приложения.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		actors:   &postgresql.ActorModel{DB: db},
		films:    &postgresql.FilmModel{DB: db},
		users:    &postgresql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на http://127.0.0.1%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn) // right or not?
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// swag init -g cmd/api/main.go

// go run ./cmd/api
