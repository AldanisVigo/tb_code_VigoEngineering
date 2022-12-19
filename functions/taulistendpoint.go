package lib

import (
	"errors"
	"fmt"
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//easyjson:json
type Categories struct {
	Categories map[int]string
}

//easyjson:json
type AddCategoryRequest struct {
	Category string
}

//export taulistendpoint
func taulistendpoint(e event.Event) uint32 {
	// Get the HTTP request
	h,err := e.HTTP()	
	if err != nil { // If we have an error getting the HTTP request
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Let the user know that we had an error
	}

	//Route the request
	err = routeRequest(h)
	if err != nil { //If there's an error while retrieving the queries
		// Set the response header's content type to application/json
		err = h.Headers().Set("Content-Type","application/json")

		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	}

	// Successful execution
	return 0;
}


// Route the request
func routeRequest(h event.HttpEvent) error {
	// Get the queries from the http event
	queries := h.Query()
	
	// Get the endpoint key value
	endpoint,err := queries.Get("endpoint")
	if err != nil {
		return err
	}

	// If the length of the endpoint param is 0
	if len(endpoint) == 0 {
		// Return a new error letting the user know what happened
		return errors.New("You must include an endpoint query parameter with your request.")
	}

	//Route to different funcs based on the selected endpoint ghetto routing
	switch endpoint { 
		case "categories":
			err = getCategories(h)
			if err != nil {
				return err
			}

			return nil
		case "addcategory": 
			err = addCategory(h)
			if err != nil {
				return err
			}
			
			h.Write([]byte(`{"added" : "true"}`))

			return nil
		default:
			_,err = h.Write([]byte("{ \"error\" : \"Invalid endpoint requested.\"}"))
			if err != nil { 
				return err
			}
			
			return nil
	}
}

func addCategory(h event.HttpEvent) error {
	//Get a reference to the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	//Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	//Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	//Create an empty incoming category request
	incomingCategoryRequest := &AddCategoryRequest{}

	//Fill it with the unmarshalled json version of the body data
	err = incomingCategoryRequest.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}
	
	// Get the categories from the database
	currentCats,err := db.Get("categories")
	if err != nil {
		return err
	}

	// Retrieve the existing list of categories
	cats := &Categories{
        Categories : map[int]string{
		},
    }

	err = cats.UnmarshalJSON(currentCats)
	if err != nil {
		return err
	}

	// Add the new category at the next available key value
	cats.Categories[len(cats.Categories) + 1] = incomingCategoryRequest.Category

	// Convert the list back to json
	j,err := cats.MarshalJSON()
	if err != nil {
		return err
	}

	// Write the list back to the database
	err = db.Put("categories",j)
	if err != nil {
		return err
	}

	// Close the databse
	err = db.Close()
	if err != nil {
		return err
	}

	// Return the json back to the user
	h.Write(j)

	// Return nil for error
	return nil
}


func getCategories(h event.HttpEvent) error {
	//Get the test database
	db, err := database.New("taulistdb")
	if err != nil { //If we encounter an error getting the database
		return err //Return the error
	}

	//Get the user JSON from the the database
	data, err := db.Get("categories")
	if err != nil { //If we encounter an error getting the current user
		return err //Return an error
	}
	
	//Close the db
	err = db.Close()
	if err != nil { //If we encounter an error while closing the database
		return err //Return the error
	}
	
	//Return a response to the caller
	w,err := h.Write([]byte(data))
	if err != nil {
		return err
	}

	//Print the results of the write
	fmt.Print(w)

	//Execution successful, return nil for error
  	return nil
}
