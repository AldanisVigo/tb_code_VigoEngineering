package lib

import (
	"fmt"

	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//export taulistendpoint
func taulistendpoint(e event.Event) uint32 {
	//Get the HTTP request
	h,err := e.HTTP()	
	if err != nil { //If we have an error getting the HTTP request
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Let the user know that we had an error
	}

	//Set the response header's content type to application/json
	err = h.Headers().Set("Content-Type","application/json")

	//Once we have the HTTP request, retrieve and return the path to the user
	// err = retrieveRequestPath(h)
	// if err != nil { //If there's an error while retrieving the path
	// 	h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	// }

	//Once we have the HTTP request, retrieve and return the request queries
	err = retrieveQueries(h)
	if err != nil { //If there's an error while retrieving the queries
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	}

	//Successful execution
	return 0;
}


//Retrieve the request query
func retrieveQueries(h event.HttpEvent) error {
	//Get the queries from the http event
	queries := h.Query()
	
	//Get the endpoint key value
	endpoint,err := queries.Get("endpoint")
	if err != nil {
		return err
	}

	//Send the endpoint query back to the client
	_,err = h.Write([]byte(endpoint))
	if err != nil {
		return err
	}

	//Execution successful, return nil for the error
	return nil
}

//Retrieve the request path
func retrieveRequestPath(h event.HttpEvent) error {
	//Get the path from the http event
	path,err := h.Path()
	if err != nil { //If we have an issue getting the path from the HTTP request
		return err //Return the error
	}

	//Write the path back to the client
	_,err = h.Write([]byte(path))
	if err != nil {
		return err
	}
	
	//Successful execution, return nil for error
	return nil
}