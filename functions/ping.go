package lib

import (
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)
var (
	testId   = uint32(5)
	testName = "someDatabase"
	testData = map[string][]byte{}
)
//export ping
func ping(e event.Event) uint32 {
	// h, err := e.HTTP()
	// if err != nil {
	// 	return 1
	// }

	// h.Write([]byte("PONG"))

	// return 0
	db, err := database.New(testName)
	if err != nil {
		return 1
	}

	err = db.Put("value/hello", []byte("Hello, world"))
	if err != nil {
		return 1
	}

	err = db.Put("value/hello2", []byte("Hello, world"))
	if err != nil {
		return 1
	}

	keys, err := db.List("value")
	if len(keys) != 2 || err != nil {
		return 1
	}

	data, err := db.Get("value/hello")
	if err != nil {
		return 1
	}

	if string(data) != "Hello, world" {
		return 1
	}

	err = db.Delete("value/hello")
	if err != nil {
		return 1
	}

	data, err = db.Get("value/hello")
	if err == nil {
		return 1
	}

	err = db.Close()
	if err != nil {
		return 1
	}

	return 0
}
