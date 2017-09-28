package rip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	DEFAULT_USER_AGENT = "rip-1.0"
)

// New returns new client
// You can provide baseUrl or you can call Base method
// 		New("http://localhost")
// is equal to
// 		New().Base("http://localhost")
func New(baseUrl ...string) (result Client) {
	result = client{
		base:      new(url.URL),
		url:       new(url.URL),
		headers:   http.Header{},
		data:      []byte{},
		client:    &http.Client{},
		userAgent: DEFAULT_USER_AGENT,
		urlValues: url.Values{},
	}

	if len(baseUrl) > 0 {
		result = result.Base(baseUrl[0])
	}
	return result.Method(http.MethodGet)
}

// client is implementation of Client interface
type client struct {
	appendSlash bool
	base        *url.URL
	url         *url.URL
	method      string
	headers     http.Header
	data        []byte
	client      *http.Client
	userAgent   string
	urlValues   url.Values
}

// AppendSlash sets whether slash should be appended automatically
func (c client) AppendSlash(append bool) Client {
	c.appendSlash = append
	return c
}

// Base sets base api url
func (c client) Base(baseUrl string) Client {
	parsed, _ := url.Parse(baseUrl)
	c.base = parsed
	return c
}

// Clients sets http client
func (c client) Client(httpClient *http.Client) Client {
	c.client = httpClient
	return c
}

// Data sets data to be sent as body of request
func (c client) Data(data interface{}) Client {
	switch d := data.(type) {
	case []byte:
		c.data = d
	case string:
		c.data = []byte(d)
	case io.Reader:
		c.data, _ = ioutil.ReadAll(d)
	default:
		c.data, _ = json.Marshal(data)
	}

	return c
}

// Do performs request and returns Response
func (c client) Do(ctx context.Context, target ...interface{}) Response {
	response := newResponse(c)

	request := c.Request().WithContext(ctx)

	// make a http call
	if httpResponse, httpError := c.client.Do(request); httpError != nil {
		response.error = httpError
		return response
	} else {
		response.response = httpResponse
	}

	var (
		body []byte
		err  error
	)
	// read body
	if body, err = ioutil.ReadAll(response.response.Body); err != nil {
		return response
	}

	response.body = body

	defer response.response.Body.Close()

	var result Response

	result = response

	// now unmarshal from json response all passed targets
	for _, t := range target {
		result = result.Unmarshal(t)
	}

	return result
}

// Header sets header to client
func (c client) Header(key, value string) Client {
	nh := copyHeaders(c.headers)
	nh.Set(key, value)

	c.headers = nh
	return c
}

// Method sets http method and adds path parts
func (c client) Method(method string, parts ...interface{}) Client {
	result := c.Path(parts...).(client)
	result.method = method

	return result
}

// Path adds path
func (c client) Path(parts ...interface{}) Client {
	strParts := []string{}

	for _, part := range parts {
		strPart := strings.TrimSpace(fmt.Sprintf("%v", part))
		if strPart == "" {
			continue
		}
		strParts = append(strParts, strPart)
	}

	joined := strings.TrimSpace(strings.Join(strParts, "/"))
	relative, _ := url.Parse(joined)

	// Should we do this?
	c.base.Path = strings.TrimRight(c.base.Path, "/") + "/"
	c.url = relative

	return c
}

// QueryParams sets url query values
func (c client) QueryValues(values url.Values) (result Client) {

	if values == nil {
		c.urlValues = url.Values{}
		return c
	}

	// first copy values
	c.urlValues = copyValues(c.urlValues)

	// assign all values
	for key, value := range values {
		for _, val := range value {
			c.urlValues.Add(key, val)
		}
	}

	return c
}

// Request returns prepared request
func (c client) Request() *http.Request {
	request, _ := http.NewRequest(c.method, c.URL().String(), bytes.NewReader(c.data))
	request.Header = copyHeaders(c.headers)
	request.Header.Set("User-Agent", c.userAgent)
	return request
}

// URL returns actual url
func (c client) URL() *url.URL {
	path := c.url.Path

	// Should we append slash?
	if c.appendSlash && path != "" {
		c.url.Path = strings.TrimRight(c.url.Path, "/") + "/"
	}

	result := c.base.ResolveReference(c.url)
	result.RawQuery = c.urlValues.Encode()

	return result
}

// UserAgent overrides default user agent
func (c client) UserAgent(agent string) Client {
	c.userAgent = agent
	return c
}

// Delete http method
func (c client) Delete(parts ...interface{}) Client {
	return c.Path(parts...).Method(http.MethodDelete)
}

// Get http method
func (c client) Get(parts ...interface{}) Client {
	return c.Method(http.MethodGet, parts...)
}

// Head http method
func (c client) Head(parts ...interface{}) Client {
	return c.Method(http.MethodHead, parts...)
}

// Options http method
func (c client) Options(parts ...interface{}) Client {
	return c.Method(http.MethodOptions, parts...)
}

// Patch http method
func (c client) Patch(parts ...interface{}) Client {
	return c.Method(http.MethodPatch, parts...)
}

// Put http method
func (c client) Put(parts ...interface{}) Client {
	return c.Method(http.MethodPut, parts...)
}

// Post http method
func (c client) Post(parts ...interface{}) Client {
	return c.Method(http.MethodPost, parts...)
}
