package http_handlers

import (
	"fmt"
	"net/http"
)

// errror is used to write error in json
// if error, when marshaling appears, handles and logs it
func responseErrorJson(rw http.ResponseWriter, statusCode int, message string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(rw, `
	{
		"error": "%s"
	}
	`, message)
}
