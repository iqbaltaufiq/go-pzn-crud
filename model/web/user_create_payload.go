package web

// a struct representing the incoming request
// when creating a new entry to db
type UserCreatePayload struct {
	Name       string `json:"name" validate:"required,min=1,max=200"`
	Occupation string `json:"occupation" validate:"required,min=1,max=200"`
}
