package login

type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type Response struct {
	Data   LoginResponse `json:"data"`
	Status int           `json:"status"`
	Error  string        `json:"error"`
}
