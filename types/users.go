package types

type User struct {
	Id       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=25"`
	Ph       string `json:"ph" validate:"required"`
}

type UserBalance struct {
	Balance int `json:"balance"`
}
