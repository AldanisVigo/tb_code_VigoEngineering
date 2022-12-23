package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/errno"
	"bitbucket.org/taubyte/go-sdk/event"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -all ${GOFILE}

//easyjson:json
type Categories struct {
	Categories []string
}

//easyjson:json
type AddCategoryRequest struct {
	Category string
}

//easyjson:json
type AdsRequest struct {
	State    string
	City     string
	Category string
}

//easyjson:json
type AddAdRequest struct {
	AdData Ad
}

//easyjson:json
type Ad struct {
	Title       string
	Description string
	PosterID    string
	City        string
	State       string
	Category    string
}

//easyjson:json
type Ads struct {
	Ads []Ad
}

//export taulistendpoint
func taulistendpoint(e event.Event) uint32 {
	// Get the HTTP request
	h, err := e.HTTP()
	if err != nil { // If we have an error getting the HTTP request
		h.Write([]byte(fmt.Sprintf("ERROR: %s\n", err))) // Let the user know that we had an error
	}

	// Route the request
	err = routeRequest(h)
	if err != nil { // If there's an error while retrieving the queries
		// Set the response header's content type to application/json
		h.Headers().Set("Content-Type", "application/json")

		// Write a response
		h.Write([]byte(fmt.Sprintf(`{ "Error" : "Error while routing your request: %s\n" }`, err))) // Send an error back to the client
	}

	// Successful execution
	return 0
}

// Route the request
func routeRequest(h event.HttpEvent) error {
	// Get the queries from the http event
	queries := h.Query()

	// Get the endpoint key value
	endpoint, err := queries.Get("endpoint")
	if err != nil {
		return err
	}

	// If the length of the endpoint param is 0
	if len(endpoint) == 0 {
		// Return a new error letting the user know what happened
		return errors.New("You must include an endpoint query parameter with your request.")
	}

	// Set he content type

	// Route to different funcs based on the selected endpoint ghetto routing
	switch endpoint {
	case "categories":
		err = getCategories(h)
		if err != nil {
			return err
		}
	case "addcategory":
		err = addCategory(h)
		if err != nil {
			return err
		}
	case "removecategory":
		err = removeCategory(h)
		if err != nil {
			return err
		}
	case "updatecategory":
		err = updateCategory(h)
		if err != nil {
			return err
		}
	case "ads":
		err = getAds(h)
		if err != nil {
			return err
		}
	case "postad":
		err = addAd(h)
		if err != nil {
			return err
		}
	case "resetcategories":
		err = resetCategories(h)
		if err != nil {
			return err
		}
	default:
		_, err = h.Write([]byte("{ \"error\" : \"Invalid endpoint requested.\"}"))
		if err != nil {
			return err
		}
	}

	// Execution sucessful
	return nil
}

func addAd(h event.HttpEvent) error {
	// Get the Body of the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	// Create an instance of an empty ad
	incomingAd := &Ad{}

	// Unmarshal the incoming bodyData into the ad object
	err = incomingAd.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	// Open the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Use the database to pull the data for the city and state and category
	existingAdsData, err := db.Get("ads/" + incomingAd.State + "/" + incomingAd.City + "/" + incomingAd.Category)

	// Unmarshal the existing ads into an Ads object
	existingAds := &Ads{
		Ads: []Ad{},
	}

	// Unmarshal the existing ads into an ads
	existingAds.UnmarshalJSON(existingAdsData)

	// Add the new ad to the existing list of ads
	existingAds.Ads = append(existingAds.Ads, *incomingAd)

	// Serialize the exising ads back into json
	existingAdsJson, err := existingAds.MarshalJSON()
	if err != nil {
		return err
	}

	// Write the existing ads back to the database
	err = db.Put("ads/"+incomingAd.State+"/"+incomingAd.City+"/"+incomingAd.Category, existingAdsJson)

	// Execution sucessful
	return nil
}

func getAds(h event.HttpEvent) error {
	// Get a reference to the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	// Create an empty incoming ads request
	incomingAdsRequest := &AdsRequest{
		City:     "",
		State:    "",
		Category: "",
	}

	// Unmarshal the incoming body data int a AdsRequest
	err = incomingAdsRequest.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	// Get the ads at the current city and state from the database
	adsForCityAndState, err := db.Get("ads/" + incomingAdsRequest.State + "/" + incomingAdsRequest.City + "/" + incomingAdsRequest.Category)
	if err != nil {
		if strings.Contains(err.Error(), errno.ErrorDatabaseKeyNotFound.String()) {
			// If the key was not found, that means there's not ads for this state and city
			h.Write([]byte(`{ "Ads" : [] }`)) // Return an empty array to the client
			return nil
		} else { // Otherwise
			return err // Return the error
		}
	}

	// If the length of the data from the database is greater than 0 bytes
	if len(adsForCityAndState) != 0 {
		// Write the data back to the client that requested it
		_, err = h.Write(adsForCityAndState)
		if err != nil {
			return err
		}
	} else { // Otherwise
		// Create an empty ads object
		adsObj := &AdsRequest{
			Ads: &[]Ad{},
		}

		// Deserialize the empty ads object into json
		adsJson, err := adsObj.MarshalJSON()
		if err != nil {
			return err
		}

		// Put the empty object json data into the database at the requested key
		err = db.Put("ads/"+incomingAdsRequest.State+"/"+incomingAdsRequest.City+"/"+incomingAdsRequest.Category, adsJson)
		if err != nil {
			return err
		}

		// Write the ads json back to the client
		_, err := h.Write(adsJson)
		if err != nil {
			return err
		}
	}

	// Execution successful
	return nil
}

func sliceContains(s *[]string, v string) (bool, error) {
	if len(v) == 0 {
		return false, errors.New("Please provide a value to check for in the slice.")
	}

	exists := false
	for _, val := range *s {
		if val == v {
			exists = true
		}
	}
	return exists, nil
}

//easyjson:json
type DeleteCategoryRequest struct {
	Index int
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func removeCategory(h event.HttpEvent) error {
	// Grab the body of the request from the http event
	body := h.Body()

	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// Sereialize the body data into a DeleteCategoryRequest
	deleteRequest := &DeleteCategoryRequest{}
	err = deleteRequest.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	// Grab the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Grab the json data for categories from the database
	categoriesJson, err := db.Get("categories")
	if err != nil {
		return err
	}

	// Seriealize the json data into a Categories object
	categoriesObj := &Categories{}

	err = categoriesObj.UnmarshalJSON(categoriesJson)
	if err != nil {
		return err
	}

	// Remove the item at the given index
	categoriesObj.Categories = RemoveIndex(categoriesObj.Categories, deleteRequest.Index)

	// Deserialize the object back into json
	updatedCategoriesJson, err := categoriesObj.MarshalJSON()
	if err != nil {
		return err
	}

	// Save the json into the database for the categories key
	err = db.Put("categories", updatedCategoriesJson)
	if err != nil {
		return err
	}

	// Close the database
	err = db.Close()
	if err != nil {
		return err
	}

	// Execution successful
	return nil
}

//easyjson:json
type CategoryUpdateRequest struct {
	Index    int
	Category string
}

func updateCategory(h event.HttpEvent) error {
	// Get the body of the http event
	body := h.Body()

	// Get the data from the body of the event
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// Seialize the body's data into a CategoryUpdateRequest object
	updateRequest := &CategoryUpdateRequest{}
	err = updateRequest.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	// Get the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Get the categories json from the database
	categoriesJson, err := db.Get("categories")
	if err != nil {
		return err
	}

	// Serialize the categories json into a Categories object
	categoriesObj := &Categories{
		Categories: []string{},
	}
	err = categoriesObj.UnmarshalJSON(categoriesJson)
	if err != nil {
		return err
	}

	// Update the value in the categoriesObj at the index
	categoriesObj.Categories[updateRequest.Index] = updateRequest.Category

	// Deserialize the categoriesObj back to json
	updatedCategories, err := categoriesObj.MarshalJSON()
	if err != nil {
		return err
	}

	if len(updatedCategories) > 0 {
		// Save the the updated categories json back to the database
		err = db.Put("categories", updatedCategories)
		if err != nil {
			return err
		}

		// Return the updated list back to the client
		_, err = h.Write(updatedCategories)
		if err != nil {
			return err
		}
	}

	// Execution successful
	return nil
}

func addCategory(h event.HttpEvent) error {
	// Get a reference to the database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Get the Body in the HTTP object
	body := h.Body()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// Close the body
	err = body.Close()
	if err != nil {
		return err
	}

	// Create an empty incoming category request
	incomingCategoryRequest := &AddCategoryRequest{}

	// Fill it with the unmarshalled json version of the body data
	err = incomingCategoryRequest.UnmarshalJSON(bodyData)
	if err != nil {
		return err
	}

	// Get the categories from the database
	currentcategories, err := db.Get("categories")
	if err != nil {
		if strings.Contains(err.Error(), errno.ErrorDatabaseKeyNotFound.String()) {
			// If the key was not found, that means there's not ads for this state and city
			// Ignore this error just keep trucking
		} else {
			return err
		}
	}

	// Create an empty categories object
	categories := &Categories{
		Categories: []string{},
	}

	// Deserialize the data from the database into the categories object
	err = categories.UnmarshalJSON(currentcategories)
	if err != nil {
		return err
	}

	// Check if the category they want to add already exist in the data structure
	exists, err := sliceContains(&categories.Categories, incomingCategoryRequest.Category)
	if err != nil {
		return err
	}

	if !exists {
		// Add the new category at the next available key value
		categories.Categories = append(categories.Categories, incomingCategoryRequest.Category)
	} else { // If the category already exists
		// Return an error to the client letting them know the category already exists.
		return errors.New("The category you are attempting to add already exists.")
	}

	// Convert the list back to json
	j, err := categories.MarshalJSON()
	if err != nil {
		return err
	}

	// Write the list back to the database
	err = db.Put("categories", j)
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

func resetCategories(h event.HttpEvent) error {
	// Ge the test database
	db, err := database.New("taulistdb")
	if err != nil {
		return err
	}

	// Delete the categories
	err = db.Delete("categories")
	if err != nil {
		return err
	}

	// Write the reset response back to the client
	_, err = h.Write([]byte(`{ "reset" : "true" }`))
	if err != nil {
		return err
	}

	return nil
}

func getCategories(h event.HttpEvent) error {
	// Get the test database
	db, err := database.New("taulistdb")
	if err != nil { // If we encounter an error getting the database
		return err // Return the error
	}

	// Get the user JSON from the the database
	data, err := db.Get("categories")
	if err != nil { // If we encounter an error getting the current user
		if strings.Contains(err.Error(), errno.ErrorDatabaseKeyNotFound.String()) {
			// If the key was not found, that means there's not ads for this state and city
			// Ignore this error
			_, err := h.Write([]byte(`{"Categories" : [], "Error" : "Categories does not exist, creating now"}`))
			if err != nil {
				return err
			}

			//Create the default categories if the categories do not exist in the db yet
			emptyCats := &Categories{
				Categories: []string{
					"Cars",
					"Apartments",
					"Books",
					"For Sale",
					"Free",
				},
			}

			//Serialize them into json
			emptyCatsJson, err := emptyCats.MarshalJSON()
			if err != nil {
				return err
			}

			//Put them in the db
			err = db.Put("categories", emptyCatsJson)
			if err != nil {
				return err
			}
		} else {
			return err // Return an error
		}
	}

	// Close the db
	err = db.Close()
	if err != nil { // If we encounter an error while closing the database
		return err // Return the error
	}

	if len(data) > 0 {
		// Return a response to the caller
		_, err := h.Write([]byte(data))
		if err != nil {
			return err
		}
	} else {
		_, err := h.Write([]byte(`{"Categories" : [] }`))
		if err != nil {
			return err
		}
	}

	// Execution successful, return nil for error
	return nil
}
