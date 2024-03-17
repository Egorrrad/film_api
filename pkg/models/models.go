package models

import "errors"

var ErrNoRecord = errors.New("подходящей записи не найдено")

type Actor struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	Birthday string   `json:"birthday"` //?????
	Films    []string `json:"films"`
}

type Film struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Rating      int      `json:"rating"`
	Actors      []string `json:"actors"`
}

type Actors struct {
	Sort   string  `json:"sort"`
	Count  int     `json:"count"`
	Actors []Actor `json:"actors"`
}

type Films struct {
	Sort  string `json:"sort"`
	Count int    `json:"count"`
	Films []Film `json:"films"`
}

type User struct {
	Role    string `json:"role"`
	Api_key string `json:"api_Key"`
}
