package db

import (
	"github.com/boltdb/bolt"
	"strconv"
)

const mapping = "mapping"
const tracks = "tracks"

type Persistent struct {
	*bolt.DB
}

func NewPersistent(db *bolt.DB) *Persistent {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(mapping))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(tracks))
		if err != nil {
			return err
		}
		return nil
	})
	return &Persistent{
		DB: db,
	}
}

func (p *Persistent) Save (key, target string) {
	p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(mapping))
		return b.Put([]byte(key), []byte(target))
	})
}
func (p *Persistent) Find(key string) (string, bool) {
	var val string
	p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(mapping))
		val = string(b.Get([]byte(key)))
		return nil
	})
	return val, val != ""
}

func (p *Persistent) Track(key string) {
	p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tracks))
		val := b.Get([]byte(key))
		num, err := strconv.Atoi(string(val))
		if err != nil {
			return err
		}
		b.Put([]byte(key), []byte(strconv.Itoa(num+1)))
		return nil
	})
}

func (p *Persistent) Visits(key string) int64 {
	var count int64
	p.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tracks))
		val := b.Get([]byte(key))
		num, err := strconv.Atoi(string(val))
		if err != nil {
			return err
		}
		count = int64(num)
		return nil
	})
	return count
}
