package models

type UserDatabase struct {
	Id, Email    string
	PasswordHash []byte
}
