package main

import "time"

type Phone struct {
	Id        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Number    string `json:"number,omitempty"`
	ContactId int64  `json:"contact_id"`
}

type Contact struct {
	Id        int64      `json:"id"`
	Firstname string     `json:"firstname,omitempty"`
	Lastname  string     `json:"lastname,omitempty"`
	Company   string     `json:"company,omitempty"`
	Address1  string     `json:"address1,omitempty"`
	Address2  string     `json:"address2,omitempty"`
	Zipcode   string     `json:"zipcode,omitempty"`
	City      string     `json:"city,omitempty"`
	Country   string     `json:"country,omitempty"`
	Email     string     `json:"email,omitempty"`
	Web       string     `json:"web,omitempty"`
	Birthday  *time.Time `json:"birthday,omitempty"`
	Tag       string     `json:"tag,omitempty"`
	Phones    []*Phone   `json:"phones,omitempty"`
}
