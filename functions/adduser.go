package lib

import (
	"fmt"
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//easyjson:json
type User struct {
	UUID string
	name string
	lname string
	age int32
}

//export adduser
func adduser(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		fmt.Print(err)
		return 1
	} 

	err = doSomeWork(h)
	if err != nil{
		h.Write([]byte(fmt.Sprintf("Add user failed with %s",err)))
		return 1
	}

	//Execution succeeded
	return 0
}

func doSomeWork(h event.HttpEvent) error {
	// //Get a reference to the database
	db, err := database.New("testdb")
	if err != nil {
		return err
	}

	//Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	//Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	//Create an empty user
	incomingUser := &User{}

	//Fill it with the unmarshalled json version of the body data
	err = incomingUser.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	//Save the user JSON to the the database
	//Ignoring errors from db.Put, h.Write, and UnmarshallJSON
	err = db.Put(incomingUser.UUID,bodyData)
	if err != nil {
		return err
	}
	
	//Close the db
	err = db.Close()
	if err != nil {
		return err
	}
	
	//Return a response to the caller
	w, err := h.Write([]byte("{ UUID : " + incomingUser.UUID + ", ADDED: true}"))
	if err != nil{
		return err
	}

	//Print out result
	fmt.Println(w)

  	return nil
}

