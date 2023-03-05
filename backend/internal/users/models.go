package users

import "time"

type User struct {
	UserDTO
	HashedPassword string
}

type UserDTO struct {
	Id        string    `json:id`
	Username  string    `json:"username"`
	PfpURL    string    `json:"pfp_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}
