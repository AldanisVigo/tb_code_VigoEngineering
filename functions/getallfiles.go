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
type FilesRequest struct {
	UUID string
	name string
	//UUID/storage/...files
}

//easyjson:json
type FilesResponse struct {
	UUID string
	name string
	files []storage.File
}

//export getallfiles
func getallfiles(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	//Get the Body in the HTTP object
	body := h.Body()
	allFilesRequestBody, err := ioutil.ReadAll(body)
	if err != nil {
		return 1
	}

	//Close the body
	err = body.Close()
	if err != nil {
		h.Write([]byte("{ Error : " + err.Error() + "}"))
		return 1
	}

	//Create an empty user
	incomingAllFilesRequest := &FilesRequest{}

	//Fill it with the unmarshalled json version of the body data
	incomingAllFilesRequest.UnmarshalJSON(allFilesRequestBody)
	
	//Get the storage for path
	filesStorage, err := storage.Get(incomingAllFilesRequest.UUID + "/" + incomingAllFilesRequest.name)
	if err != nil {
		h.Write([]byte("{ UUID : " + incomingAllFilesRequest.UUID + ", Error : " + err.Error() + "}"))
		return 1
	}

	//Get the files from the storage at that path
	files, err := filesStorage.ListFiles()
	if err != nil {
		h.Write([]byte(" { UUID : " + incomingAllFilesRequest.UUID + ", Error : " + err.Error() + "}"))
		return 1
	}

	//Attach the files to the response
	filesResponse := &FilesResponse{
		UUID: incomingAllFilesRequest.UUID,
		name : incomingAllFilesRequest.name,
		files : files,
	}
	
	//Get the serialized json from the response we created
	filesResponseJson, err := filesResponse.MarshalJSON()
	if err != nil {
		h.Write([]byte(" { UUID : " + incomingAllFilesRequest.UUID + ", Error : " + err.Error() + "}"))
		return 1
	}
	
	//Return a response to the caller
	w,err := h.Write([]byte(filesResponseJson))
	if err != nil {
		return 1
	}

	//Print results of calling Write
	fmt.Print(w)
	
  	return 0
}