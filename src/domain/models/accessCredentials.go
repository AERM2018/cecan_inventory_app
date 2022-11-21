package models

type AccessCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}
