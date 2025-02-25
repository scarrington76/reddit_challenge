package model

// User is the model used to track a user. If a db was implemented, the sql tags would be used
type User struct {

	// Name is the name of the user
	Name string `sql:"name"`

	// Posts is the current # of posts by the user
	Posts int `sql:"posts"`
}
