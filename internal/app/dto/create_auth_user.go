package dto

type CreateAuthUserInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm,omitempty"`
	Name            string `json:"name"`
}

type CreateAuthUserOutput struct {
	PublicID  string `json:"public_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
