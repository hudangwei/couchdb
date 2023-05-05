package couchdb

type DatabaseService interface {
	AllDocs(params *QueryParameters) (*ViewResponse, error)
	AllDesignDocs() ([]DesignDocument, error)
	Find(*FindArgs, interface{}) error
	CreateIndex(*Index) (*CouchIndexBody, error)
	Rev(id string) (string, error)
	Get(doc CouchDoc, id string) error
	Put(doc CouchDoc) (*DocumentResponse, error)
	Post(doc CouchDoc) (*DocumentResponse, error)
	Delete(doc CouchDoc) (*DocumentResponse, error)
	Store(doc CouchDoc) (*DocumentResponse, error)
	MultiStore(docs []CouchDoc) error
	PutAttachmentToDoc(doc CouchDoc, path string) (*DocumentResponse, error)
	Bulk(docs []CouchDoc) ([]DocumentResponse, error)
	Purge(req map[string][]string) (*PurgeResponse, error)
	View(name string) ViewService
	Seed([]DesignDocument) error
	IndependAttachment(docid, name, rev string) (*IndependAttachment, error)
	IndependAttachmentMeta(docid, name, rev string) (*IndependAttachment, error)
	PutIndependAttachment(docid string, att *IndependAttachment, rev string) (*DocumentResponse, error)
	DeleteIndependAttachment(docid, name, rev string) (*DocumentResponse, error)
}

// ViewService is an interface for dealing with a view inside a CouchDB database.
type ViewService interface {
	Get(name string, params QueryParameters) (*ViewResponse, error)
	Post(name string, keys []string, params QueryParameters) (*ViewResponse, error)
}
