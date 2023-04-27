package couchdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	fileNameMap    = "map.js"
	fileNameReduce = "reduce.js"
)

// Parse takes a location and parses all design documents with corresponding views.
// The folder structure must look like this.
//
//   design
//   |-- player
//   |   |-- byAge
//   |   |   |-- map.js
//   |   |   `-- reduce.js
//   |   `-- byName
//   |       `-- map.js
//   `-- user
//       |-- byEmail
//       |   |-- map.js
//       |   `-- reduce.js
//       `-- byUsername
//           `-- map.js
func (c *CouchDBClient) Parse(dirname string) ([]DesignDocument, error) {
	docs := []DesignDocument{}
	// get all directories inside location which will become separate design documents
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		designDocumentName := dir.Name()
		ff, err := ioutil.ReadDir(filepath.Join(dirname, designDocumentName))
		if err != nil {
			return nil, err
		}
		d := DesignDocument{
			Document: Document{
				ID: fmt.Sprintf("_design/%s", designDocumentName),
			},
			Language: langJavaScript,
			Views:    map[string]DesignDocumentView{},
		}
		for _, j := range ff {
			viewName := j.Name()
			// create new view inside design document
			view := DesignDocumentView{}
			// get map function
			pathMap := filepath.Join(dirname, designDocumentName, viewName, fileNameMap)
			bMap, err := ioutil.ReadFile(pathMap)
			if err != nil {
				return nil, err
			}
			view.Map = string(bMap)
			// get reduce function only if it exists
			pathReduce := filepath.Join(dirname, designDocumentName, viewName, fileNameReduce)
			if _, err := os.Stat(pathReduce); err != nil {
				// ignore error that file does not exist but return other errors
				if !os.IsNotExist(err) {
					return nil, err
				}
			} else {
				bReduce, err := ioutil.ReadFile(pathReduce)
				if err != nil {
					return nil, err
				}
				view.Reduce = string(bReduce)
			}
			d.Views[viewName] = view
		}
		docs = append(docs, d)
	}
	return docs, nil
}
