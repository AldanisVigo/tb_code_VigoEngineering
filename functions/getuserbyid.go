package lib

import (
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//easyjson:json
type UserRequest struct {
	UUID string
}

//export getuserbyid
func getuserbyid(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	// //Get a reference to the database
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

	//Create an empty user
	incomingUser := &UserRequest{}

	//Fill it with the unmarshalled json version of the body data
	incomingUser.UnmarshalJSON(bodyData)

	//Save the user JSON to the the database
	data, err := db.Get(incomingUser.UUID)
	if err != nil {
		return 1
	}
	
	//Close the db
	err = db.Close()
	if err != nil {
		return 1
	}
	
	//Return a response to the caller
	h.Write([]byte(data))

  	return 0
}
