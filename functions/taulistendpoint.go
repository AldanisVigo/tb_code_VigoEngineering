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
type CategoriesList struct {
	categories []string
}

/*
	ModifyCategories allows you to change the categories in the CategoriesList struct
*/
func (cl *CategoriesList) ModifyCategories(newCategories []string){
	cl.categories = newCategories
}

//export taulistendpoint
func taulistendpoint(e event.Event) uint32 {
	// Get the HTTP request
	h,err := e.HTTP()	
	if err != nil { // If we have an error getting the HTTP request
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Let the user know that we had an error
	}

	// Set the response header's content type to application/json
	err = h.Headers().Set("Content-Type","application/json")

	// Once we have the HTTP request, retrieve and return the path to the user
	// err = retrieveRequestPath(h)
	// if err != nil { //If there's an error while retrieving the path
	// 	h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	// }

	// Once we have the HTTP request, retrieve and return the request queries
	err = retrieveQueryParams(h)
	if err != nil { //If there's an error while retrieving the queries
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	}

	// Successful execution
	return 0;
}


// Retrieve the params from the request query
func retrieveQueryParams(h event.HttpEvent) error {
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

	switch endpoint { 
		case "categories":
			err = getCategories(&h)
			if err != nil {
				return err
			}

			return nil
		case "addcategory": 
			err = addCategory(&h)
			if err != nil {
				return err
			}

			return nil
		default:
			_,err = h.Write([]byte("{ \"error\" : \"Invalid endpoint requested.\"}"))
			if err != nil { 
				return err
			}
			
			return nil
	}
}

//easyjson:json
type Categories struct {
	categories map[int]string
}

//easyjson:json
type AddCategoryRequest struct {
	category  string
}
//export addCategory
func addCategory(h *event.HttpEvent) error {
	// Open the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Retrieve the new category from the request
	requestBody := h.Body()
	requestBodyData,err := ioutil.ReadAll(requestBody)
	if err != nil {
		return err
	}
	req := &AddCategoryRequest{}
	req.UnmarshalJSON(requestBodyData)
	newCat := req.category

	// Get the categories from the database
	currentCats,err := db.Get("categories")
	if err != nil {
		return err
	}

	// Retrieve the existing list of categories
	cats := &Categories{
        // categories : map[string]string{
        //     "Ti1": "hello",
        // },
    }
	err = cats.UnmarshalJSON(currentCats)
	if err != nil {
		return err
	}

	// Add the new category at the next available key value
	cats.categories[len(cats.categories) + 1] = newCat

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

//export getCategories
func getCategories(h *event.HttpEvent) error {
	//Get the test database
	db, err := database.New("testdb")
	if err != nil { //If we encounter an error getting the database
		return err //Return the error
	}

	//Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body) //Read the contents of the request body
	if err != nil { //If we encounter an error reading the contents of the request body
		return err //Return the error
	}

	//Close the body
	err = body.Close() 
	if err != nil { //If we encounter an error closing the request body
		return err //Return the error
	}

	//Create an empty user
	incomingUserRequest := &UserRequest{}

	//Fill it with the unmarshalled json version of the body data
	incomingUserRequest.UnmarshalJSON(bodyData)

	//Get the user JSON from the the database
	data, err := db.Get(incomingUserRequest.UUID)
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
