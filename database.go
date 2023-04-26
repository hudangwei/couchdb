package couchdb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

type Database struct {
	Client *CouchDBClient
	Host   string
	Name   string
}

// AllDesignDocs returns all design documents from database.
// http://stackoverflow.com/questions/2814352/get-all-design-documents-in-couchdb
func (db *Database) AllDesignDocs() ([]DesignDocument, error) {
	startKey := fmt.Sprintf("%q", "_design/")
	endKey := fmt.Sprintf("%q", "_design0")
	includeDocs := true
	q := QueryParameters{
		StartKey:    &startKey,
		EndKey:      &endKey,
		IncludeDocs: &includeDocs,
	}
	res, err := db.AllDocs(&q)
	if err != nil {
		return nil, err
	}
	docs := make([]interface{}, len(res.Rows))
	for index, row := range res.Rows {
		docs[index] = row.Doc
	}
	designDocs := make([]DesignDocument, len(docs))
	b, err := json.Marshal(docs)
	if err != nil {
		return nil, err
	}
	return designDocs, json.Unmarshal(b, &designDocs)
}

// AllDocs returns all documents in selected database.
// http://docs.couchdb.org/en/latest/api/database/bulk-api.html
func (db *Database) AllDocs(params *QueryParameters) (*ViewResponse, error) {
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s/%s/_all_docs?%s", db.Host, url.PathEscape(db.Name), q.Encode())
	response := &ViewResponse{}
	err = db.Client.Get(u, response)
	return response, err
}

func (db *Database) Get(doc CouchDoc, id string) error {
	u := fmt.Sprintf("%s/%s/%s", db.Host, url.PathEscape(db.Name), url.PathEscape(id))
	return db.Client.Get(u, doc)
}

func (db *Database) Put(doc CouchDoc) (*DocumentResponse, error) {
	u := fmt.Sprintf("%s/%s/%s", db.Host, url.PathEscape(db.Name), url.PathEscape(doc.GetID()))
	response := &DocumentResponse{}
	err := db.Client.Put(u, doc, response)
	return response, err
}

func (db *Database) Post(doc CouchDoc) (*DocumentResponse, error) {
	u := fmt.Sprintf("%s/%s", db.Host, url.PathEscape(db.Name))
	response := &DocumentResponse{}
	err := db.Client.Post(u, doc, response)
	return response, err
}

func (db *Database) Delete(doc CouchDoc) (*DocumentResponse, error) {
	u := fmt.Sprintf("%s/%s/%s?rev=%s", db.Host, url.PathEscape(db.Name), url.PathEscape(doc.GetID()), doc.GetRev())
	response := &DocumentResponse{}
	err := db.Client.Delete(u, response)
	return response, err
}

// Bulk allows to create and update multiple documents
// at the same time within a single request. The basic operation is similar to
// creating or updating a single document, except that you batch
// the document structure and information.
func (db *Database) Bulk(docs []CouchDoc) ([]DocumentResponse, error) {
	bulk := BulkDoc{
		Docs: docs,
	}
	u := fmt.Sprintf("%s/%s/_bulk_docs", db.Host, url.PathEscape(db.Name))
	response := []DocumentResponse{}
	err := db.Client.Post(u, bulk, &response)
	return response, err
}

// Purge permanently removes the references to deleted documents from the database.
// http://docs.couchdb.org/en/1.6.1/api/database/misc.html
func (db *Database) Purge(req map[string][]string) (*PurgeResponse, error) {
	u := fmt.Sprintf("%s/%s/_purge", db.Host, url.PathEscape(db.Name))
	response := &PurgeResponse{}
	err := db.Client.Post(u, req, response)
	return response, err
}

func (db *Database) View(name string) ViewService {
	u := fmt.Sprintf("%s/%s/_design/%s/", db.Host, url.PathEscape(db.Name), url.PathEscape(name))
	return &View{
		URL:    u,
		Client: db.Client,
	}
}

// Seed makes sure all your design documents are up to date.
func (db *Database) Seed(cache []DesignDocument) error {
	// query all docs to get all design documents
	designDocs, err := db.AllDesignDocs()
	if err != nil {
		return err
	}
	difference := diff(cache, designDocs)
	// remove all deletions
	for _, doc := range difference.deletions {
		if _, err := db.Delete(&doc); err != nil {
			return err
		}
	}
	// update all changes
	for _, doc := range difference.changes {
		// get design document first to get current revision
		var old DesignDocument
		if err := db.Get(&old, doc.ID); err != nil {
			return err
		}
		// update document with new version
		doc.Rev = old.Rev
		if _, err := db.Put(&doc); err != nil {
			return err
		}
	}
	// add all additions
	for _, doc := range difference.additions {
		if _, err := db.Put(&doc); err != nil {
			return err
		}
	}
	return nil
}

type difference struct {
	additions []DesignDocument
	changes   []DesignDocument
	deletions []DesignDocument
}

func diff(cache, db []DesignDocument) difference {
	di := difference{
		additions: []DesignDocument{},
		changes:   []DesignDocument{},
		deletions: []DesignDocument{},
	}
	// check for additions changes
	// design document is in cache but not in db
	for _, c := range cache {
		exists := false
		existsButDifferent := false
		for _, d := range db {
			if d.ID == c.ID {
				exists = true
				// check for different map/reduce and language
				// do not check for different revision
				if !reflect.DeepEqual(c.Views, d.Views) {
					existsButDifferent = true
				}
			}
		}
		if !exists {
			di.additions = append(di.additions, c)
		} else if existsButDifferent {
			di.changes = append(di.changes, c)
		}
	}
	// check for deletions
	// design document is in db but not in cache
	for _, d := range db {
		exists := false
		for _, c := range cache {
			if d.ID == c.ID {
				exists = true
			}
		}
		// do not delete internal design documents like _auth
		if !exists && !strings.HasPrefix(d.Name(), "_") {
			di.deletions = append(di.deletions, d)
		}
	}
	return di
}

func (db *Database) Find(args *FindArgs, out interface{}) error {
	u := fmt.Sprintf("%s/%s/_find", db.Host, url.PathEscape(db.Name))
	return db.Client.Post(u, args, out)
}

func (db *Database) CreateIndex(args *Index) (*CouchIndexBody, error) {
	u := fmt.Sprintf("%s/%s/_index", db.Host, url.PathEscape(db.Name))
	response := &CouchIndexBody{}
	err := db.Client.Post(u, args, response)
	return response, err
}