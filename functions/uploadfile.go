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
type FileUploadRequest struct {
	UUID string
	filePath string
	file string
}

func (req *FileUploadRequest) ModifyUUID(uuid string) {
  req.UUID = uuid
}

func (req *FileUploadRequest) ModifyFilePath(filePath string){
	req.filePath = filePath
}

func (req *FileUploadRequest) ModifyFile(file string){
	req.file = file
}

func getFileUploadRequestFromJson(json string,req *FileUploadRequest) error {
	if len(json) == 0{
		return errors.New("The json provided is empty")
	}

	//Get the parts out of the json
	_,after,openingCharFound := strings.Cut(json,"{") 
	if !openingCharFound { //If the opening { is not found
		//Let the user know we couldn't parse the json
		return errors.New("Could not parse JSON starting json character { not found.")
	}else { //Otherwise, we can grab everything after the { character and cut it by }
		//We're interested in everything before the } character
		before,_,closingCharFound := strings.Cut(after,"}")
		if !closingCharFound { //If the } character was not found 
			//Let the user know the json could not be parsed because of the missing character
			return errors.New("Could not parse JSON closing json character } not found.")
		}

		//If we are able to get everything between the { ... }, we can then split by commas
		//Create a map of key value pairs that is empty
		var mapOfKeyValuePairs = make(map[string]string)

		//Split the ... content by commas to separate the key:value pairs
		keyValuePairsWithColonSep := strings.Split(before, ",")

		//Iterate through the key:value pairs
		for _ ,value := range keyValuePairsWithColonSep {
			//Split them by the colon to separate the key and the value
			keyVal := strings.Split(value, ":")
			
			//Isolate the key
			key := strings.Split(string(keyVal[0]), "\"")[1]
			
			//Isolate the value
			val := strings.Split(string(keyVal[1]), "\"")[1]

			//Modify the key in the map of key value pairs
			mapOfKeyValuePairs[key] = val
		}

		//Modify the values in the empty req
		req.ModifyFile(mapOfKeyValuePairs["file"])
		req.ModifyFilePath(mapOfKeyValuePairs["filePath"])
		req.ModifyUUID(mapOfKeyValuePairs["UUID"])

		//Return a nil error
		return nil
	}
}

//export uploadfile
func uploadfile(e event.Event) uint32 {
	//Get the http object from the event
  	h, err := e.HTTP()
	if err != nil {	//If there's an error
		//Print the error
		fmt.Println(err.Error())
		return 1 //Eject
	}

	//Upload the file in the request
	err = uploadFileInRequest(h)
	if err != nil { //If there's any errors while uploading the file to storage
		//Write a response to the client with the error
		h.Write([]byte("{ \"error\" : \"" + err.Error() + "\"}"))
	}

	//Execution successful
	return 0
}

func uploadFileInRequest(h event.HttpEvent) error {
 	//Get a reference to our existing storage bucket
	testStorage, err := storage.New("teststorage")
	if err != nil { //If there's an error
		return err //Return the error
	}

	//Get the Body in the HTTP object
	body := h.Body()
	
	//Read the contents of the body of the request
	fileUploadRequestContents, err := ioutil.ReadAll(body)
	if err != nil { //If there's an error while reading
		return err //Return the error
	}

	err = body.Close() //Close the body
	if err != nil { //If there's an error while closing the body
		return err //Return the error
	}

	//Create an empty FileUploadRequest
	req := &FileUploadRequest{}

	//Fill it up with stuff from the fileUploadRequestContents
	err = getFileUploadRequestFromJson(string(fileUploadRequestContents),req)
	if err != nil { //If we get an error getting the UploadRequest from the json
		return err //Return the error
	}
	
	//Save the file in the json request to the file storage at the uuid/name/file path
	file := testStorage.File(req.UUID + "/" + req.filePath)
	
	//Add the file and get the version of the file
	version , err := file.Add([]byte(req.file), true)
	if err != nil { //If there's an error while adding the file to the storage
		return err //Return the error
	}

	//Print the version of the file
	fmt.Print(version)
	
	//Return a response to the caller
	w, err := h.Write([]byte("{ \"UUID\" : \"" + req.UUID + "\" , \"filename\" : \"" + req.filePath + "\",\"file_uploaded\" : true }"))
	if err != nil { //If there's an error while writing a response back 
		return err //Return the error
	}

	//Print the result of writing a response
	fmt.Print(w)

	//Successful operation, return nil for error
  	return nil
}