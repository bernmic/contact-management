package model

type Contact struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
