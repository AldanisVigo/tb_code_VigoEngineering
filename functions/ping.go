package lib

import (
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//export ping
func ping(e event.Event) uint32 {
	db, err := database.New("testdb")
	
	if err != nil {
		return 1
	}

	// err = db.Put("value/hello", []byte("Hello, world"))
	// if err != nil {
	// 	return 1
	// }

	// err = db.Put("value/hello2", []byte("Hello, world"))
	// if err != nil {
	// 	return 1
	// }

	keys, err := db.List("value")
	if len(keys) != 2 || err != nil {
		return 1
	}

	//Get HTTP from the event
	h, err := e.HTTP()
	if err != nil { //If we get an err 
		return 1 //roll out
	}

	//Otherwise write the keys to the page cuz why not?
	// h.Write([]byte(strings.Join(keys, ",")))
	data, err := db.Get("value/hello")
	if err != nil {
		return 1
	}

	h.Write(data)
	// if string(data) != "Hello, world" {
	// 	return 1
	// }

	// err = db.Delete("value/hello")
	// if err != nil {
	// 	return 1
	// }


	// data, err := db.Get("value/hello")
	// if err == nil {
	// 	return 1
	// }
	// h.Write([]byte(data))
	// err = db.Close()
	// if err != nil {
	// 	return 1
	// }

	return 0
}
