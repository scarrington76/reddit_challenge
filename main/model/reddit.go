package model

// NewResponse is the struct of the response when hitting the reddit API for listings.
// This struct is used only for unmarshaling responses.
type NewResponse struct {
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Title  string `json:"title"`
				Ups    int    `json:"ups"`
				ID     string `json:"id"`
				Author string `json:"author"`
			} `json:"data"`
		} `json:"children"`
		Before interface{} `json:"before"`
	} `json:"data"`
}
