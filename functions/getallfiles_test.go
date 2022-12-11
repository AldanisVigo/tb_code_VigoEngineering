package lib

// import (
// 	"testing"

// 	eventSym "bitbucket.org/taubyte/go-sdk-symbols/event"
// 	"bitbucket.org/taubyte/go-sdk/common"
// )

// func TestGetAllFiles(t *testing.T){
// 	//Mock the get all files event
// 	mockEventCall := eventSym.MockData{
// 		EventId:    0,
// 		EventType:  common.EventTypeHttp,
// 		Body: 		[]byte(`{
// 			UUID : 'aldanisvigo'
// 		}`),
// 		Headers:    map[string]string{},
// 		Queries:    map[string]string{},
// 		Host:       "",
// 		Method:     "POST",
// 		Path:       "/getallfiles",
// 		ReturnCode: 0,
// 	}.Mock()

// 	//Mock a call to getallfiles func
// 	if getallfiles(0) == 0 {
// 		t.Errorf("Error: %s",string(mockEventCall.ReturnBody))
// 		return
// 	}

// 	//Get a string of the expected response
// 	expectedResponse := string([]byte(`{
// 		FILES : [
// 			something,
// 			something_else
// 		]
// 	}`))

// 	//Check if the return body of the mocked event call matches the expected response
// 	if string(mockEventCall.ReturnBody) == expectedResponse { //If they don't match
// 		//Generate an error message
// 		t.Errorf("Got: %s expected %s",string(mockEventCall.ReturnBody), expectedResponse)
// 		return
// 	}
// }
	