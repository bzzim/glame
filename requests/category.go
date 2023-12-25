package requests

type Category struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
}
