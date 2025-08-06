package dto

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
