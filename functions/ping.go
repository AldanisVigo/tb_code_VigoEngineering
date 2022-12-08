package lib

import (
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//export ping
func ping(e event.Event) uint32 {
	//Get the database reference
	db, err := database.New("testdb")
	if err != nil {
		return 1
	}

	//Get the data from the database
	data,err := db.Get("value/hello")
	if err != nil {
		return 1
	}

	//Get HTTP from the event
	h, err := e.HTTP()
	if err != nil { //If we get an err 
		return 1 //roll out
	}

	//Send the data back to the browser
	h.Write(data)
	return 0
}
