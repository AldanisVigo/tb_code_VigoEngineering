package lib

import (
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

type User struct {
    Name  string
	Last string
	Age int32
}

//export wok
func wok(e event.Event) uint32 {
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	db, err := database.New("testdb")
	if err != nil {
		return 1
	}

	data, err := db.Get("value/hello")
	if err != nil {
		return 1
	}

 	h.Write([]byte(data))
  
  	return 0
}