package lib

import (
	"fmt"
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/event"
	"bitbucket.org/taubyte/go-sdk/storage"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//easyjson:json
type FileUploadRequest struct {
	UUID string
	name string
	file []byte

	//UUID/filename/file
}

//export uploadfile
func uploadfile(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	// //Get a reference to our existing storage bucket
	testStorage, err := storage.New("teststorage")
	if err != nil {
		return 1
	}

	//Get the Body in the HTTP object
	body := h.Body()
	fileUploadRequestContents, err := ioutil.ReadAll(body)
	if err != nil {
		return 1
	}

	//Close the body
	err = body.Close()
	if err != nil {
		return 1
	}

	//Create an empty user
	incomingFileUploadRequest := &FileUploadRequest{}

	//Fill it with the unmarshalled json version of the body data
	incomingFileUploadRequest.UnmarshalJSON(fileUploadRequestContents)

	//Save the file in the json request to the file storage at the uuid/name/file path
	file := testStorage.File(incomingFileUploadRequest.UUID + "/" + incomingFileUploadRequest.name)
	version , err := file.Add(incomingFileUploadRequest.file, true)
	if err != nil {
		return 1
	}

	//Print the version of the file
	fmt.Print(version)
	
	//Return a response to the caller
	h.Write([]byte("{ UUID : " + incomingFileUploadRequest.UUID + ", FILE_NAME: " + incomingFileUploadRequest.name + ",FILE_UPLOADED: true}"))
	
  	return 0
}