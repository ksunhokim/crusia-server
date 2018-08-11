package store

import "time"

type User struct {
	ID       int
	Username string
	Email    string
	Passhash string
}

type SaveData struct {
	UserID  int
	Edited  time.Time
	Payload string
}

type Store interface {
	GetUser(id int) (*User, error)
	GetUserByUsername(name string) (*User, error)
	CreateUser(*User) (*User, error)
	UpdateUser(*User) error
	DeleteUser(*User) error

	GetSaveData(userid int) (*SaveData, error)
	CreateSaveData(*SaveData) (*SaveData, error)
	UpdateSaveData(*SaveData) error
	DeleteSaveData(*SaveData) error
}
