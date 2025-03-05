package model

// Post is the model used to track a post's votes. If a db was implemented, the sql tags would be used
type Post struct {

	// Name is the name of the post
	Name string `sql:"name"`

	// Votes is the number of current upvotes that have been tracked
	Votes int `sql:"votes"`
}
