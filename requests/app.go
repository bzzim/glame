package requests

type App struct {
	Name        string `json:"name" form:"name"`
	Url         string `json:"url" form:"url"`
	Icon        string `json:"icon,omitempty" form:"-"`
	Description string `json:"description" form:"description"`
	IsPublic    bool   `json:"isPublic" form:"isPublic"`
}

type AppOrder struct {
	Apps []ItemOrder `json:"apps"`
}
