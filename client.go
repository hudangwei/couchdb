package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type CouchDBClient struct {
	client *http.Client
	host   string
	user   string
	pwd    string
}

func NewClient(host, user, pwd string) *CouchDBClient {
	return &CouchDBClient{http.DefaultClient, host, user, pwd}
}

func (c *CouchDBClient) Use(name string) DatabaseService {
	return &Database{
		Host:   c.host,
		Name:   name,
		Client: c,
	}
}

func (c *CouchDBClient) Head(rawurl string) (*http.Response, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("HEAD", uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.pwd)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CouchDBClient) GetRaw(rawurl string) (*http.Response, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.pwd)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CouchDBClient) Get(rawurl string, out interface{}) error {
	return c.do(rawurl, "GET", nil, out)
}

func (c *CouchDBClient) Post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

func (c *CouchDBClient) Put(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PUT", in, out)
}

func (c *CouchDBClient) PutWithData(rawurl string, data io.Reader, out interface{}, contentType string) error {
	return c.doWithoutEncode(rawurl, "PUT", data, out, contentType)
}

func (c *CouchDBClient) Delete(rawurl string, out interface{}) error {
	return c.do(rawurl, "DELETE", nil, out)
}

func (c *CouchDBClient) do(rawurl, method string, in, out interface{}) error {
	body, err := c.open(rawurl, method, in)
	if err != nil {
		return err
	}
	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

func (c *CouchDBClient) doWithoutEncode(rawurl, method string, in io.Reader, out interface{}, contentType string) error {
	body, err := c.openWithoutEncode(rawurl, method, in, contentType)
	if err != nil {
		return err
	}
	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

func (c *CouchDBClient) doWithoutDecode(rawurl, method string, in interface{}) ([]byte, error) {
	body, err := c.open(rawurl, method, in)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	out, _ := ioutil.ReadAll(body)
	return out, nil
}

func (c *CouchDBClient) open(rawurl, method string, in interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.pwd)
	if in != nil {
		decoded, derr := json.Marshal(in)
		if derr != nil {
			return nil, derr
		}
		buf := bytes.NewBuffer(decoded)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(len(decoded))
		req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}

func (c *CouchDBClient) openWithoutEncode(rawurl, method string, in io.Reader, contentType string) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), in)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.pwd)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}

func handleRsp(rsp *http.Response, err error) ([]byte, error) {
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
