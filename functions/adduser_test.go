package lib

import (
	"testing"

	databaseSym "bitbucket.org/taubyte/go-sdk-symbols/database"
	eventSym "bitbucket.org/taubyte/go-sdk-symbols/event"
	"bitbucket.org/taubyte/go-sdk/common"
)


func TestAddUser(t *testing.T){
	//Mock the event
	mockEventCall := eventSym.MockData{
		EventId:    0,
		EventType:  common.EventTypeHttp,
		Body:       []byte(`{
			"name": "Aldanis",
			"lname": "Vigo",
			"age": 32,
			"UUID": "aldanisvigo"
		}`),
		Headers:    map[string]string{},
		Queries:    map[string]string{},
		Host:       "",
		Method:     "POST",
		Path:       "/adduser",
		ReturnCode: 0,
	}.Mock()
	
	//Mock the database
	databaseSym.Mock(0,"testdb_fail",map[string][]byte{})

	//Mock call adduser with the event id 0
	if adduser(0) == 0 { //If we get an error
		//Generate an error response
		t.Errorf("Error Expected: %s",string(mockEventCall.ReturnBody))
		return 
	}

	//Mock the database with test id 0 which matches the id of our mocked event
	databaseSym.Mock(0,"testdb",map[string][]byte{})

	//If we use event with id 0 and we get 1 instead
	if adduser(0) == 1 {
	 	//Generate an error 
		t.Errorf("Failed because: %s",string(mockEventCall.ReturnBody))
		return
	}

	//Generate an expected response 
	expectedResponse := "{ UUID : aldanisvigo, ADDED : true }"

	//If the mocked event call returns the expected response
	if string(mockEventCall.ReturnBody) == expectedResponse {
		t.Errorf("Got: %s expected %s",string(mockEventCall.ReturnBody), expectedResponse)
		return 
	}
}