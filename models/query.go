package models

type Queries struct {
	Queries []Query `json:"queries"`
}

type Query struct {
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Template string `json:"template"`
}
