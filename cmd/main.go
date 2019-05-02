package main

import (
	"fmt"
	"short-url/server"
)

func main() {
	// uncomment the following lines for a persistent data store (BoltDB).
	// Note: I'm setting a global variable which is very bad. Again, just a quick and dirty implementation.
	// There is a known bug with the `Visits` and the persistent store. It always returns 0.

	//d, err := bolt.Open("shortner.db", 0600, nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer d.Close()
	//server.MainDB = db.NewPersistent(d)

	fmt.Println("listening on", server.BaseURL)
	if e := server.ListenAndServe(); e != nil {
		fmt.Println("ERROR:", e)
	}
}
