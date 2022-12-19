package lib

import (
	"errors"
	"fmt"

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

// func addNewCategory(h *event.HttpEvent) error {
// 	// Get a reference to the database
// 	db, err := database.New("taulistdb")
// 	if err != nil {
// 		return err
// 	}

// 	// Get the Body in the HTTP object
// 	body := h.Body()
// 	bodyData, err := ioutil.ReadAll(body)
// 	if err != nil {
// 		return err
// 	}

// 	// Close the body
// 	err = body.Close()
// 	if err != nil {
// 		return err
// 	}

// 	// Create an empty array of strings to hold the categories
// 	addCategoryRequest := &AddCategoryRequest{}

// 	// Unmarshal the request 
// 	// err = addCategoryRequest.UnmarshalJSON(bodyData)
// 	err = json.Unmarshal(bodyData,addCategoryRequest)
// 	if err != nil {
// 		return err
// 	}

// 	// Get the current list of categories
// 	catsJson,err := db.Get("categories")
// 	if err != nil {
// 		return err
// 	}

// 	//Create an empty slice of strings
// 	catsList := &CategoriesList{}

// 	//Fill it with the unmarshaled json from the database
// 	// err = catsList.UnmarshalJSON(catsJson)
// 	err = json.Unmarshal(catsJson,catsList)
// 	if err != nil {
// 		return err
// 	}

// 	// Add the requested category to the list of existing categories
// 	catsList.categories = append(catsList.categories, addCategoryRequest.category)

// 	// Get the json string for the categories
// 	catsListJson,err := catsList.MarshalJSON()
// 	if err != nil {
// 		return err
// 	}

// 	//Save the user JSON to the the database
// 	//Ignoring errors from db.Put, h.Write, and UnmarshallJSON
// 	err = db.Put("categories",catsListJson)
// 	if err != nil {
// 		return err
// 	}
	
// 	//Close the db
// 	err = db.Close()
// 	if err != nil {
// 		return err
// 	}
	
// 	//Return a response to the caller
// 	w, err := h.Write([]byte("{ category : \"" + addCategoryRequest.category + "\", \"added\" : \"true\" }"))
// 	if err != nil{
// 		return err
// 	}

// 	//Print out result
// 	fmt.Println(w)

// 	//Execution successful, return nil for error
//   	return nil
// }

