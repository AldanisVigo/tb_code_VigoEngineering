package lib

import (
	"encoding/json"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

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

	j,err := json.MarshalIndent(data,""," ")
	if err != nil{
		return 1
	}

 	h.Write([]byte(j))
  
  	return 0
}