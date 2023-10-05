package models

type AuthDataAuthorized struct {
	IsAuth     bool    `json:"isAuth"`
	Login      string  `json:"login"`
	Email      string  `json:"email"`
	AccessType string  `json:"accessType"`
}

type AuthDataUnauthorized struct {
	IsAuth bool `json:"isAuth"`
}
