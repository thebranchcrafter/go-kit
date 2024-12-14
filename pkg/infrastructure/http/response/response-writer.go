package http_response

import "net/http"

// ResponseWriter is the interface that abstracts the response writing behavior.
type ResponseWriter interface {
	WriteErrorResponse(w http.ResponseWriter, err error, httpStatus int, previousError error)
	WriteResponse(w http.ResponseWriter, payload interface{}, httpStatus int)
}
