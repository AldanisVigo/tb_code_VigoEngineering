package lib

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/errno"
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

	// Get the database 
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}


	switch endpoint { 
		case "categories":
			// Retrieve the categories from the database
			cats,err := retrieveCategories(&db)
			if err != nil { // If there's an error retrieving the categories from the database
				return err // Return the error
			}

			// Send the categories back to the client
			_,err = h.Write([]byte(cats))
			if err != nil { // If there's an error sending the categories to the client
				return err // Return the error
			}
			
			// Execution succeeded, return nil for error
			return nil
		case "addcategory":
			// Retrieve the query parameter for the category
			cat,err := queries.Get("category")
			if err != nil { // If there's an error pulling the query parameter for the category
				return err // Return the error
			}else{ // Otherwise
				// Add the category to the category list
				err = addCategory(&db,cat,&h)
				if err != nil { // if there's an error while adding the category to the category list
					return err // Return the error
				}
				
				_,err = h.Write([]byte("Added category"))
			}

			// Execution succeeded, return nil for error
			return nil
		default:
			// Send an empty json object back to the client
			_,err = h.Write([]byte("{ \"error\" : \"Invalid endpoint requested.\"}"))
			if err != nil { // If there's an error writing the json back to the client
				return err // Return the error
			}
			
			// Execution succeeded, return nil for error
			return nil
	}
}


/*
	Given a json string, parse the json string and set the fields in the passed
	CategoriesList Object
*/
func serializeCategoriesJson(json string,catList *CategoriesList,h *event.HttpEvent) error {
	if len(json) == 0 { // If the length of the provided json is 0
		// Return an error letting the user know that their json is empty
		return errors.New("Error serializing the categories json to a CategoriesList instance: The json provided was empty.")
	}

	// Serialize the json into the catList object
	err := catList.UnmarshalJSON([]byte(json))
	if err != nil {
		return err
	}

	//Operation successful return nil for error
	return nil
}

// Retrieve all categories stored in the taulist database
func retrieveCategories(db *database.Database) (string, error) {
	// Get the json data in the categories
	cats, err := db.Get("categories")
	if err != nil { //If there's an error
		// Get the missing key error from errno
		var missingKeyError errno.Error = errno.ErrorDatabaseKeyNotFound
	
		// If the generated error contains the ErrorDatabaseKeyNotFound error
		if strings.Contains(err.Error(), missingKeyError.String()) {
			return "{ \"categories\" : [] }", nil // Return an empty json object and nil for the error
		}
		
		// For all other errors do the same thing for now
		// TODO: Add more robust error handling later
		return "{ \"categories\" : [] }", nil //Return an empty json object and nil for the error
	}

	if len(cats) == 0 { // If there's no cats
		// Return an empty json object, and nil for the error
		return "{ \"categories\" : [] }", nil
	}

	// Execution successful, return nil for the error
	return string(cats), nil
}

// Add a category to the taulist databse
func addCategory(db *database.Database, category string,h *event.HttpEvent) error {
	// Retrieve the vales of the categories
	currentCats,err := retrieveCategories(db)

	// Create an empty categories list
	catListObj := &CategoriesList{}

	// TODO: Serialize the current categories into a CategoriesList object
	err = serializeCategoriesJson(currentCats,catListObj,h)
	if err != nil {
		return err
	}

	// h.Write([]byte(catListObj.categories));
	for _,val := range catListObj.categories {
		_,err = h.Write([]byte(val))
		if err != nil {
			return err
		}
	}

	// Put the value in the categories 
	err = db.Put("categories",[]byte(currentCats))
	if err != nil {
		return err
	}

	// Execution successful, return nil for the error
	return nil
}

// Retrieve the request path
func retrieveRequestPath(h event.HttpEvent) error {
	// Get the path from the http event
	path,err := h.Path()
	if err != nil { // If we have an issue getting the path from the HTTP request
		return err // Return the error
	}

	// Write the path back to the client
	_,err = h.Write([]byte(path))
	if err != nil {
		return err
	}
	
	// Successful execution, return nil for error
	return nil
}
