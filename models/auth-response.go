package models

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ResponseSuccess struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type ResponseError struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}
