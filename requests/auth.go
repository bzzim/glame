package requests

type Login struct {
	Duration string `json:"duration" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Validate struct {
	Token string `json:"token" binding:"required"`
}
