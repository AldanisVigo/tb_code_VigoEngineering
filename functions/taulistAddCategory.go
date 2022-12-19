package lib

import (
	"io/ioutil"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

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