package web

type HttpResponse struct {
	Code   int
	Status string
	Data   interface{}
}
