package lib

import (
	"io/ioutil"
	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//export adduser
func adduser(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	//Get a reference to the database
	db, err := database.New("testdb")
	if err != nil {
		return 1
	}

	//Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return 1
	}

	//Close the body
	err = body.Close()
	if err != nil {
		return 1
	}



	//Close the db
	err = db.Close()
	if err != nil {
		return 1
	}

	//Return what we get
 	h.Write([]byte(bodyData))
  
  	return 0
}
