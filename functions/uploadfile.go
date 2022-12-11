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
	testStorage, err := storage.New("teststorage")
	if err != nil { //If there's an error
		//Write a response to the client with the error
		h.Write([]byte("{ \"error\": \"" + err.Error() + "\", \"msg\" : \"There was an rerror opening the test storage in the dFunction.\"}"))
		return 1 //Eject
	}

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
	incomingFileUploadRequest := &FileUploadRequest{}

	//Fill it with the unmarshalled json version of the body data
	err = incomingFileUploadRequest.UnmarshalJSON(fileUploadRequestContents)
	if err != nil { //If there's an error while serializing the JSON into a FileUploadRequest
		//Send a response back to the client letting them know about the error
		h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error serializing your request into a FileUploadRequest in te dFunction.\"}")) 
		return 1 //Eject
	}

	w,err := h.Write(fileUploadRequestContents)
	fmt.Print(w)

	// //Save the file in the json request to the file storage at the uuid/name/file path
	// file := testStorage.File(incomingFileUploadRequest.UUID + "/" + incomingFileUploadRequest.name)
	
	// //Add the file and get the version of the file
	// version , err := file.Add(incomingFileUploadRequest.file, true)
	// if err != nil { //If there's an error while adding the file to the storage
	// 	//Write a response to the client letting them know about the error
	// 	h.Write([]byte("{ \"error\" : \"" + err.Error() + "\", \"msg\" : \"There was an error adding the file to the test storage in the dFunction.\"}"))
	// 	return 1 //Eject
	// }

	// //Print the version of the file
	// fmt.Print(version)
	
	// //Return a response to the caller
	// w, err := h.Write([]byte("{ \"UUID\" : \"" + incomingFileUploadRequest.UUID + "\" , \"filename\" : \"" + incomingFileUploadRequest.name + "\",\"file_uploaded\" : true }"))
	// if err != nil { //If there's an error while writing a response back 
	// 	fmt.Print(err) //Print the error
	// }

	// //Print the result of writing a response
	// fmt.Print(w)

	//Successful operation
  	return 0
}