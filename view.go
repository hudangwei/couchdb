package couchdb

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// View performs actions and certain view documents
type View struct {
	URL    string
	Client *CouchDBClient
}

// Get executes specified view function from specified design document.
func (v *View) Get(name string, params QueryParameters) (*ViewResponse, error) {
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s_view/%s?%s", v.URL, name, q.Encode())
	response := &ViewResponse{}
	err = v.Client.Get(uri, response)
	return response, err
}

// Post executes specified view function from specified design document.
// Unlike View.Get for accessing views, View.Post supports
// the specification of explicit keys to be retrieved from the view results.
func (v *View) Post(name string, keys []string, params QueryParameters) (*ViewResponse, error) {
	content := struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	}
	// create query string
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s_view/%s?%s", v.URL, name, q.Encode())
	response := &ViewResponse{}
	err = v.Client.Post(url, content, response)
	return response, err
}
