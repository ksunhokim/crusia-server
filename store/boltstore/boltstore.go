package boltstore

import (
	"github.com/asdine/storm"
	"github.com/sdbx/crusia-server/store"
)

type BoltStore struct {
	db *storm.DB
}

func New(db *storm.DB) *BoltStore {
	return &BoltStore{
		db: db,
	}
}

func (b *BoltStore) Init() error {
	err := b.db.Init(&store.User{})
	if err != nil {
		return err
	}

	return b.db.Init(&store.SaveData{})
}

func (b *BoltStore) GetUser(id int) (*store.User, error) {
	var u store.User
	err := b.db.One("ID", id, &u)
	return &u, err
}

func (b *BoltStore) GetUserByUsername(name string) (*store.User, error) {
	var u store.User
	err := b.db.One("Username", name, &u)
	return &u, err
}

func (b *BoltStore) CreateUser(u *store.User) error {
	u.ID = 0
	return b.db.Save(u)
}

func (b *BoltStore) UpdateUser(u *store.User) error {
	return b.db.Update(u)
}

func (b *BoltStore) DeleteUser(u *store.User) error {
	return b.db.DeleteStruct(u)
}

func (b *BoltStore) GetSaveData(userid int) (*store.SaveData, error) {
	var d store.SaveData
	err := b.db.One("UserID", userid, &d)
	return &d, err
}

func (b *BoltStore) CreateSaveData(d *store.SaveData) error {
	return b.db.Save(d)
}

func (b *BoltStore) UpdateSaveData(d *store.SaveData) error {
	return b.db.Update(d)
}

func (b *BoltStore) DeleteSaveData(d *store.SaveData) error {
	return b.db.DeleteStruct(d)
}
