package dto

type CreateAuthUserInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Name            string `json:"name"`
}

type CreateAuthUserOutput struct {
	ID       int64  `json:"id"`
	PublicID string `json:"public_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
