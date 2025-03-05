package posts

// response is the dto used to publish a post's votes.
type response struct {

	// Name is the name of the post
	Name string `json:"name"`

	// Votes is the number of current upvotes that have been tracked
	Votes int `json:"votes"`
}
