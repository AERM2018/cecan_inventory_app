package models

type AuthClaims struct {
	Id       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}
