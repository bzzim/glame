package requests

type Bookmark struct {
	Name       string `json:"name" form:"name"`
	Url        string `json:"url" form:"url"`
	Icon       string `json:"icon,omitempty"`
	CategoryId int    `json:"categoryId" form:"categoryId"`
	IsPublic   bool   `json:"isPublic" form:"isPublic"`
}

type BookmarkOrder struct {
	Bookmarks []ItemOrder `json:"bookmarks"`
}
