package user

import "server/internal/ws/client"

type User struct {
	ConnectionClient *client.Client

	Username string
}

func NewUser (username string, client *client.Client) *User{
	return &User{Username: username, ConnectionClient: client}
}