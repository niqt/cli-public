package httpx

import "net/http"

type Method interface {
	String() string
}

type httpMethod string

func (m httpMethod) String() string {
	return string(m)
}

const (
	MethodGet    httpMethod = "GET"
	MethodPost   httpMethod = "POST"
	MethodPut    httpMethod = "PUT"
	MethodDelete httpMethod = "DELETE"
)

type Service interface {
	Method() Method
	Path() string
	Handler() http.HandlerFunc
}
