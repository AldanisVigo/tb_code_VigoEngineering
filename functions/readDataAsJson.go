package lib

//Import the necessary libraries
import (
	"encoding/json"
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//export ping
func readDataAsJson(e event.Event) uint32 {
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

	//Generate a json response	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return 1
	}

	//Send the data back to the browser
	h.Write(jsonData)

	//Return 0 cuz IDK, just do it
	return 0
}
