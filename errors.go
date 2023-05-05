package couchdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents API-level errors, reported by CouchDB as
//
//	{"error": <ErrorCode>, "reason": <Reason>}
type Error struct {
	Method     string // HTTP method of the request
	URL        string // HTTP URL of the request
	StatusCode int    // HTTP status code of the response

	// These two fields will be empty for HEAD requests.
	ErrorCode string // Error reason provided by CouchDB
	Reason    string // Error message provided by CouchDB
}

func (e *Error) Error() string {
	if e.ErrorCode == "" {
		return fmt.Sprintf("%v %v: %v", e.Method, e.URL, e.StatusCode)
	}
	return fmt.Sprintf("%v %v: (%v) %v: %v",
		e.Method, e.URL, e.StatusCode, e.ErrorCode, e.Reason)
}

// NotFound checks whether the given errors is a DatabaseError
// with StatusCode == 404. This is useful for conditional creation
// of databases and documents.
func NotFound(err error) bool {
	return ErrorStatus(err, http.StatusNotFound)
}

// Unauthorized checks whether the given error is a DatabaseError
// with StatusCode == 401.
func Unauthorized(err error) bool {
	return ErrorStatus(err, http.StatusUnauthorized)
}

// Conflict checks whether the given error is a DatabaseError
// with StatusCode == 409.
func Conflict(err error) bool {
	return ErrorStatus(err, http.StatusConflict)
}

// ErrorStatus checks whether the given error is a DatabaseError
// with a matching statusCode.
func ErrorStatus(err error, statusCode int) bool {
	dberr, ok := err.(*Error)
	return ok && dberr.StatusCode == statusCode
}

func parseError(req *http.Request, resp *http.Response) error {
	var reply struct{ Error, Reason string }
	if req.Method != "HEAD" {
		if err := readBody(resp, &reply); err != nil {
			return fmt.Errorf("couldn't decode CouchDB error: %v", err)
		}
	}
	return &Error{
		Method:     req.Method,
		URL:        req.URL.String(),
		StatusCode: resp.StatusCode,
		ErrorCode:  reply.Error,
		Reason:     reply.Reason,
	}
}

func readBody(resp *http.Response, v interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		resp.Body.Close()
		return err
	}
	return resp.Body.Close()
}
