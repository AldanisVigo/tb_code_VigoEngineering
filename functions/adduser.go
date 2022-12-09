package lib

import (
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
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

	
	incomingUser := &User{}
	incomingUser.UnmarshalJSON(bodyData)

	db.Put(incomingUser.UUID,bodyData)

	h.Write([]byte("{ UUID : " + incomingUser.UUID + ", ADDED: true}"))

	// //Close the db
	// err = db.Close()
	// if err != nil {
	// 	return 1
	// }

	//Return what we get
 	// h.Write([]byte(incomingUser.name + " " + incomingUser.lname + " - Age: " + string(incomingUser.age)))
  
  	return 0
}

func ( User ) MarshalJSON() ([]byte, error) { return nil, nil }
func (* User ) UnmarshalJSON([]byte) error { return nil }
func ( User ) MarshalEasyJSON(w *jwriter.Writer) {}
func (* User ) UnmarshalEasyJSON(l *jlexer.Lexer) {}

type EasyJSON_exporter_User *User