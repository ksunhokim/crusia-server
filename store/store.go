package store

import "time"

type User struct {
	ID       int    `storm:"id,increment"`
	Username string `storm:"index"`
	Email    string `storm:"index"`
	Passhash string
}

type SaveData struct {
	UserID  int `storm:"id"`
	Edited  time.Time
	Payload string
}

type Store interface {
	GetUser(id int) (*User, error)
	GetUserByUsername(name string) (*User, error)
	CreateUser(*User) error
	UpdateUser(*User) error
	DeleteUser(*User) error

	GetSaveData(userid int) (*SaveData, error)
	CreateSaveData(*SaveData) error
	UpdateSaveData(*SaveData) error
	DeleteSaveData(*SaveData) error
}
