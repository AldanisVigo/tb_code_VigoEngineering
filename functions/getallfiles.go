package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

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

func (req *FilesRequest) ModifyUUID(uuid string) error {
	if len(uuid) == 0 {
		return errors.New("The provided UUID is empty")
	}

	req.UUID = uuid
	return nil
}

func (req *FilesRequest) ModifyName(name string) error {
	if len(name) == 0 {
		return errors.New("The provided name is empty")
	}

	req.name = name
	return nil
}

//easyjson:json
type FilesResponse struct {
	UUID string
	name string
	files string
}

func getFileRequestFromJson(json string, req *FilesRequest) error{
	//If the incoming json is empty
	if len(json) == 0 {
		//Return an error
		return errors.New("The json provided is empty.")
	}

	_,after,containsStarting := strings.Cut(json,"{")
	if !containsStarting {
		return errors.New("Error parsing json. Missing starting { json character.")
	} else {
		before,_,containsEnding := strings.Cut(after,"}")
		if !containsEnding {
			return errors.New("Error parsing json. Missing ending } json character.")
		}
		
		//Make a map to hold the key value pairs
		keyValPairsMap := make(map[string]string)

		//Split the keyValPairs by ,
		keyValPairs := strings.Split(before, ",")
		
		//Iterate through them
		for _,value := range keyValPairs {
			kvp := strings.Split(value, ":") //Split it by the :
			key := strings.Split(kvp[0], "\"")[1] //Grab the key
			val := strings.Split(kvp[1], "\"")[1] //Grab the value
			keyValPairsMap[key] = val //Save it to the map
		}

		//Modify the request sender's UUID
		req.ModifyUUID(keyValPairsMap["UUID"])

		//And the file's name
		req.ModifyName(keyValPairsMap["name"])

		//There were no errors, return nil
		return nil
	}
}

//export getallfiles
func getallfiles(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
		if err != nil {
		return 1
	}

	//Attempt to retrieve the requested files
	err = retrieveRequestedFiles(h)
	if err != nil { //If there's an error while attempting to get the files
		//Send a response to the client letting them know there was an error
		h.Write([]byte(fmt.Sprintf("{\"error\" : \"%s\"}",err.Error())))
		return 1
	}

	return 0
}

func retrieveRequestedFiles(h event.HttpEvent) error {

	//Get the Body in the HTTP object
	body := h.Body()
	allFilesRequestBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	//Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	//Create an empty file request
	filesReq := &FilesRequest{}

	//Fill it in with the data from the request body
	err = getFileRequestFromJson(string(allFilesRequestBody),filesReq)
	if err != nil {
		return err
	}
	
	//Get the test storage
	storageRef, err := storage.New("teststorage")
	if err != nil {
		return err
	}

	//Get the file
	file := storageRef.File(filesReq.UUID + "/" + filesReq.name)
	storageFile, err := file.GetFile()
	if err != nil {
		return err
	}

	//Read the file into a byte array
	fileContents, err := ioutil.ReadAll(storageFile)
	if err != nil {
		return err
	}

	//Set the response headers content type to application/json
	h.Headers().Set("Content-Type","application/json")
	
	//Return a response to the caller
	w,err := h.Write([]byte("{ \"UUID\" : \"" + filesReq.UUID + "\", \"name\" : \"" + filesReq.name + "\", \"file\" : \"" + string(fileContents) + "\" }"))
	if err != nil {
		return err
	}

	// //Print results of calling Write
	fmt.Print(w)
	
	//No errors so return nil
  	return nil
}