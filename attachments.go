package couchdb

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type IndependAttachment struct {
	Name string    // Filename
	Type string    // MIME type of the Body
	MD5  []byte    // MD5 checksum of the Body
	Body io.Reader // The body itself
}

func (db *Database) IndependAttachment(docid, name, rev string) (*IndependAttachment, error) {
	if docid == "" {
		return nil, fmt.Errorf("couchdb.GetAttachment: empty docid")
	}
	if name == "" {
		return nil, fmt.Errorf("couchdb.GetAttachment: empty attachment Name")
	}
	var u string
	if rev == "" {
		u = fmt.Sprintf("%s/%s/%s/%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(name))
	} else {
		u = fmt.Sprintf("%s/%s/%s/%s?rev=%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(name), url.PathEscape(rev))
	}
	resp, err := db.Client.GetRaw(u)
	if err != nil {
		return nil, err
	}
	att, err := attFromHeaders(name, resp)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	att.Body = resp.Body
	return att, nil
}

func (db *Database) IndependAttachmentMeta(docid, name, rev string) (*IndependAttachment, error) {
	if docid == "" {
		return nil, fmt.Errorf("couchdb.GetAttachment: empty docid")
	}
	if name == "" {
		return nil, fmt.Errorf("couchdb.GetAttachment: empty attachment Name")
	}
	var u string
	if rev == "" {
		u = fmt.Sprintf("%s/%s/%s/%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(name))
	} else {
		u = fmt.Sprintf("%s/%s/%s/%s?rev=%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(name), url.PathEscape(rev))
	}
	resp, err := db.Client.Head(u)
	if err != nil {
		return nil, err
	}
	return attFromHeaders(name, resp)
}

func attFromHeaders(name string, resp *http.Response) (*IndependAttachment, error) {
	att := &IndependAttachment{Name: name, Type: resp.Header.Get("content-type")}
	md5 := resp.Header.Get("content-md5")
	if md5 != "" {
		if len(md5) < 22 || len(md5) > 24 {
			return nil, fmt.Errorf("couchdb: Content-MD5 header has invalid size %d", len(md5))
		}
		sum, err := base64.StdEncoding.DecodeString(md5)
		if err != nil {
			return nil, fmt.Errorf("couchdb: invalid base64 in Content-MD5 header: %v", err)
		}
		att.MD5 = sum
	}
	return att, nil
}

func (db *Database) PutIndependAttachment(docid string, att *IndependAttachment, rev string) (*DocumentResponse, error) {
	if docid == "" {
		return nil, fmt.Errorf("couchdb.PutAttachment: empty docid")
	}
	if att.Name == "" {
		return nil, fmt.Errorf("couchdb.PutAttachment: empty attachment Name")
	}
	if att.Body == nil {
		return nil, fmt.Errorf("couchdb.PutAttachment: nil attachment Body")
	}

	var u string
	if rev == "" {
		u = fmt.Sprintf("%s/%s/%s/%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(att.Name))
	} else {
		u = fmt.Sprintf("%s/%s/%s/%s?rev=%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(att.Name), url.PathEscape(rev))
	}

	response := &DocumentResponse{}
	err := db.Client.PutWithData(u, att.Body, response, att.Type)
	return response, err
}

func (db *Database) DeleteIndependAttachment(docid, name, rev string) (*DocumentResponse, error) {
	if docid == "" {
		return nil, fmt.Errorf("couchdb.PutAttachment: empty docid")
	}
	if name == "" {
		return nil, fmt.Errorf("couchdb.PutAttachment: empty name")
	}

	u := fmt.Sprintf("%s/%s/%s/%s?rev=%s", db.Host, url.PathEscape(db.Name), url.PathEscape(docid), url.PathEscape(name), url.PathEscape(rev))
	response := &DocumentResponse{}
	err := db.Client.Delete(u, response)
	return response, err
}
