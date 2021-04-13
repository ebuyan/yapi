package yandex

import "github.com/tidwall/buntdb"

type Store struct {
	db *buntdb.DB
}

func NewStore() Store {
	db, _ := buntdb.Open("db")
	return Store{db: db}
}

func (s Store) setToken(token string) {
	s.db.Update(func(tx *buntdb.Tx) (err error) {
		tx.Set("oauth_token", token, nil)
		return
	})
}

func (s Store) getToken() (ok bool, token string) {
	s.db.View(func(tx *buntdb.Tx) (err error) {
		token, err = tx.Get("oauth_token")
		ok = err == nil
		return
	})
	return
}
