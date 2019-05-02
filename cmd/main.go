package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"short-url/db"
	"short-url/server"
)

func main() {
	d, err := bolt.Open("shortner.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	server.MainDB = db.NewPersistent(d)

	fmt.Println("listening on", server.BaseURL)
	if e := server.ListenAndServe(); e != nil {
		fmt.Println("ERROR:", e)
	}
}
