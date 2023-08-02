package main

import "net/http"

func NOTmain() {
	// create a new handler
	handler := HttpHandler{}
	// listen and serve
	http.ListenAndServe(":9000", handler)
}

// create a handler struct
type HttpHandler struct{}

// implement `ServeHTTP` method on `HttpHandler` struct
func (h HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// create response binary data
	data := []byte("Hello World!") // slice of bytes
	// write `data` to response
	res.Write(data)
}

/*
func ListenAndServe(addr string, handler Handler) error

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

ServeHTTP(res http.ResponseWriter, req *http.Request)

type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
*/
