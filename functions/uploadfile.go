package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"bitbucket.org/taubyte/go-sdk/event"
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

func getFileUploadRequestFromJson(json string) (*FileUploadRequest,error){
	if len(json) == 0{
		return nil, errors.New("The json provided is empty")
	}

	//Get the parts out of the json
	_,after,openingCharFound := strings.Cut(json,"{") 
	if !openingCharFound { //If the opening { is not found
		//Let the user know we couldn't parse the json
		return nil, errors.New("Could not parse JSON starting json character { not found.")
	}else { //Otherwise, we can grab everything after the { character and cut it by }
		//We're interested in everything before the } character
		before,_,closingCharFound := strings.Cut(after,"}")
		if !closingCharFound { //If the } character was not found 
			//Let the user know the json could not be parsed because of the missing character
			return nil, errors.New("Could not parse JSON closing json character } not found.")
		}

		//If we are able to get everything between the { ... }, we can then split by commas
		//Create a map of key value pairs that is empty
		mapOfKeyValuePairs := make(map[string]string)
		//Split the ... content by commas to separate the key:value pairs
		keyValuePairsWithColonSep := strings.Split(before, ",")

		//Iterate through the key:value pairs
		for _ ,value := range keyValuePairsWithColonSep {
			//Split them by the colon to separate the key and the value
			keyVal := strings.Split(value, ":")
			//Isolate the key
			key := keyVal[0]
			//Isolate the value
			val := keyVal[1]
			//Add an entry on the map at the key set to the value
			mapOfKeyValuePairs[key] = val
		}

		//Make an object of type FileUploadRequest and fill it with the information in the map
		//And a nil error
		return &FileUploadRequest{
			UUID : mapOfKeyValuePairs["UUID"],
			filePath: mapOfKeyValuePairs["filePath"],
			file : mapOfKeyValuePairs["file"],
		}, nil
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

	// //Get a reference to our existing storage bucket
	// testStorage, err := storage.New("teststorage")
	// if err != nil { //If there's an error
	// 	//Write a response to the client with the error
	// 	h.Write([]byte("{ \"error\": \"" + err.Error() + "\", \"msg\" : \"There was an rerror opening the test storage in the dFunction.\"}"))
	// 	return 1 //Eject
	// }

	//Get the Body in the HTTP object
	body := h.Body()
	//Read the contents of the body of the request
	fileUploadRequestContents, err := ioutil.ReadAll(body)
	if err != nil { //If there's an error while reading
		//Write a response to the client letting them know of the error
		h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error reading the body of your request in the dFunction.\"}"))
		return 1 //Eject
	}

	err = body.Close() //Close the body
	if err != nil { //If there's an error while closing the body
		//Write a response to the client letting them know about the error
		h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error closing the body or your request in the dFunction.\"}"))
		return 1 //Eject
	}

	//Create an empty FileUploadRequest
	// incomingFileUploadRequest := &FileUploadRequest{}

	//Fill it with the unmarshalled json version of the body data
	// err = incomingFileUploadRequest.UnmarshalJSON(fileUploadRequestContents)
	// if err != nil { //If there's an error while serializing the JSON into a FileUploadRequest
	// 	//Send a response back to the client letting them know about the error
	// 	h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error serializing your request into a FileUploadRequest in te dFunction.\"}")) 
	// 	return 1 //Eject
	// }

	//Set the header's Content-Type of the response to application/json
	// err = h.Headers().Set("Content-Type","application/json")
	// if err != nil { //If there's an error setting the header's content type
	// 	return 1 //Eject
	// }

	req,err := getFileUploadRequestFromJson(string(fileUploadRequestContents))
	if err != nil {
		h.Write([]byte("{ \"error\" : \"" +  err.Error() + "\" }"))
		return 1
	}
	
	//Write the json response back to the client
	w,err := h.Write([]byte("{ \"UUID\" : \"" + req.UUID + "\", \"path\" : \"" + req.filePath +"\", \"file\" : \"" + req.file + "\"}"))
	fmt.Print(w)

	// // //Save the file in the json request to the file storage at the uuid/name/file path
	// file := testStorage.File(incomingFileUploadRequest.UUID + "/" + incomingFileUploadRequest.UUID + "/" + incomingFileUploadRequest.filePath)
	
	// // //Add the file and get the version of the file
	// version , err := file.Add([]byte(incomingFileUploadRequest.file), true)
	// if err != nil { //If there's an error while adding the file to the storage
	// 	//Write a response to the client letting them know about the error
	// 	h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error adding the file to the test storage in the dFunction.\"}"))
	// 	return 1 //Eject
	// }

	// //Print the version of the file
	// fmt.Print(version)
	
	// //Return a response to the caller
	// w, err = h.Write([]byte("{ \"UUID\" : \"" + incomingFileUploadRequest.UUID + "\" , \"filename\" : \"" + incomingFileUploadRequest.filePath + "\",\"file_uploaded\" : true }"))
	// if err != nil { //If there's an error while writing a response back 
	// 	fmt.Print(err) //Print the error
	// }

	// // //Print the result of writing a response
	// fmt.Print(w)

	//Successful operation
  	return 0
}