package models

type Themes struct {
	Themes []Theme `json:"themes"`
}

type Theme struct {
	Name     string     `json:"name"`
	IsCustom bool       `json:"isCustom"`
	Colors   ThemeColor `json:"colors"`
}

type ThemeColor struct {
	Accent     string `json:"accent"`
	Background string `json:"background"`
	Primary    string `json:"primary"`
}
