package lib

import (
	"fmt"
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

func getCategories(h *event.HttpEvent) error {
	//Get the test database
	db, err := database.New("testdb")
	if err != nil { //If we encounter an error getting the database
		return err //Return the error
	}

	//Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body) //Read the contents of the request body
	if err != nil { //If we encounter an error reading the contents of the request body
		return err //Return the error
	}

	//Close the body
	err = body.Close() 
	if err != nil { //If we encounter an error closing the request body
		return err //Return the error
	}

	//Create an empty user
	incomingUserRequest := &UserRequest{}

	//Fill it with the unmarshalled json version of the body data
	incomingUserRequest.UnmarshalJSON(bodyData)

	//Get the user JSON from the the database
	data, err := db.Get(incomingUserRequest.UUID)
	if err != nil { //If we encounter an error getting the current user
		return err //Return an error
	}
	
	//Close the db
	err = db.Close()
	if err != nil { //If we encounter an error while closing the database
		return err //Return the error
	}
	
	//Return a response to the caller
	w,err := h.Write([]byte(data))
	if err != nil {
		return err
	}

	//Print the results of the write
	fmt.Print(w)

	//Execution successful, return nil for error
  	return nil
}
