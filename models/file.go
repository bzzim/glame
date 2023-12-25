package models

type Files struct {
	Files []File `json:"files"`
}

type File struct {
	Name     string      `json:"name"`
	Msg      FileMsg     `json:"msg"`
	Paths    FilePaths   `json:"paths"`
	Template interface{} `json:"template"`
	IsJSON   bool        `json:"isJSON"`
}

type FilePaths struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

type FileMsg struct {
	Created string `json:"created"`
	Found   string `json:"found"`
}
