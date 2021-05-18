package store

import "github.com/tidwall/buntdb"

type Store struct {
	*buntdb.DB
}

func NewStore() Store {
	db, _ := buntdb.Open("db")
	return Store{db}
}

func (s Store) SetToken(token string) {
	s.Update(func(tx *buntdb.Tx) (err error) {
		tx.Set("oauth_token", token, nil)
		return
	})
}

func (s Store) GetToken() (ok bool, token string) {
	s.View(func(tx *buntdb.Tx) (err error) {
		token, err = tx.Get("oauth_token")
		ok = err == nil
		return
	})
	return
}
