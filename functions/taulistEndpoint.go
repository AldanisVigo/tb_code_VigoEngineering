package lib

import (
	"errors"
	"fmt"

	"bitbucket.org/taubyte/go-sdk/database"
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
	err = retrieveQueryParams(h)
	if err != nil { //If there's an error while retrieving the queries
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n",err))) //Send an error back to the client
	}

	//Successful execution
	return 0;
}


//Retrieve the params from the request query
func retrieveQueryParams(h event.HttpEvent) error {
	//Get the queries from the http event
	queries := h.Query()
	
	//Get the endpoint key value
	endpoint,err := queries.Get("endpoint")
	if err != nil {
		return err
	}

	//If the length of the endpoint param is 0
	if len(endpoint) == 0 {
		//Return a new error letting the user know what happened
		return errors.New("You must include an endpoint query parameter with your request.")
	}


	switch endpoint { 
		case "categories":
			//Retrieve the categories from the database
			cats,err := retrieveCategories(db)
			if err != nil { //If there's an error retrieving the categories from the database
				return err //Return the error
			}

			//Send the categories back to the client
			_,err = h.Write([]byte(cats))
			if err != nil { //If there's an error sending the categories to the client
				return err //Return the error
			}
			
			//Execution succeeded, return nil for error
			return nil
		default:
			//Send an empty json object back to the client
			_,err = h.Write([]byte("{}"))
			if err != nil { //If there's an error writing the json back to the client
				return err //Return the error
			}
			
			//Execution succeeded, return nil for error
			return nil
	}

	//Send the endpoint query back to the client
	_,err = h.Write([]byte(endpoint))
	if err != nil {
		return err
	}

	//Execution successful, return nil for the error
	return nil
}

//easyjson:json
type CategoriesList struct {
	categories []string
}

func (cl *CategoryList) ModifyCategories(newCategories []string){
	cl.categories = newCategories
}

func serializeCategoriesJson(json string,catList *CategoriesList) error {
	if len(json) == 0 { //If the length of the provided json is 0
		//Return an error letting the user know that their json is empty
		return errors.New("Error serializing the categories json to a CategoriesList instance: The json provided was empty.")
	}

	_,after,containsOpening := strings.Cut(json,"{")
	if !containsOpening {
		return errors.New("Error serializing the categories json to a CategoryList instance: The json provided is missing the opening {")
	}else{
		before,_,containsClosing := strings.Cut(after,"}")
		if !containsClosing {
			return errors.New("Error serializing the categories json to a CategoryList instance: The json provided is missing the closing }")
		}
		keyValPairsWithColon,err := strings.Split(before,":")
		
	}

}

//Retrieve all categories stored in the taulist database
func retrieveCategories(db database.Database) ([]byte, error) {
	//Get the json data in the categories
	cats, err := db.Get("categories")
	if err != nil {
		return nil, err
	}

	if len(cats) == 0 { //If there's no cats
		//Return an empty json object, and nil for the error
		return []byte("{}"), nil
	}

	//Execution successful, return nil for the error
	return cats, nil
}

//Add a category to the taulist databse
func addCategory(db database.Database, category string) error {
	//Retrieve the vales of the categories
	currentCats :=

	//Put the value in the categories 
	w, err = db.Put("categories",category)
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