package main

type Users struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Pass     string `json:"pass"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}
