package dto

type ValidateCreateAccount struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}