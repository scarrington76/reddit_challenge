package model

// NewResponse is the struct of the response when hitting the reddit API for listings.
// This struct is used only for unmarshaling reddit responses.
type NewResponse struct {

	// Kind is the kind of type
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Kind string `json:"kind"`
			// Data is the object that holds all of the information for a post/listing
			Data struct {
				// Title is the name of the post
				Title string `json:"title"`
				// ps are the number of upvotes for a post
				Ups int `json:"ups"`
				// ID is the unique identifier for the reddit post
				ID string `json:"id"`
				// Author is the username for the reddit post
				Author string `json:"author"`
			} `json:"data"`
		} `json:"children"`
		Before interface{} `json:"before"`
	} `json:"data"`
}
