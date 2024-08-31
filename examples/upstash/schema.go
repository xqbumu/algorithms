package main

import "time"

type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
}

func NewUser() User {
	return User{
		Name:           "user",
		FavoriteNumber: time.Now().UnixMilli(),
	}
}
