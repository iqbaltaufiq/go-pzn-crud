package web

type UserUpdatePayload struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1,max=200"`
}
