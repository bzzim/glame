package responses

type Login struct {
	Token string `json:"token"`
}

type ValidateToken struct {
	IsValid bool `json:"isValid"`
}

type Validate struct {
	Token ValidateToken `json:"token"`
}
