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
	DefaultUserAgent = "rip-1.0"
)

// New returns new client
// You can provide baseUrl or you can call Base method
// 		New("http://localhost")
// is equal to
// 		New().Base("http://localhost")
func New(baseUrl ...string) (result Client) {
	result = client{
		base:      new(url.URL),
		headers:   http.Header{},
		data:      []byte{},
		client: func() *http.Client {
			return &http.Client{
				Transport: &http.Transport{
					DisableKeepAlives: true,
				},
			}
		},
		userAgent: DefaultUserAgent,
		urlValues: url.Values{},
		parts:     make([]string, 0),
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
	before      func(r *http.Request)
	method      string
	headers     http.Header
	data        []byte
	client      func() *http.Client
	userAgent   string
	urlValues   url.Values
	parts       []string
}

// BeforeSend sets callback before sending request
func (c client) BeforeSend(bf func(r *http.Request)) Client {
	c.before = bf
	return c
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

// Clients sets http client function. This function is called every `Do` call so we use new client.
// This is needed in environment where we share client across goroutines
func (c client) Client(httpClientFunc func() *http.Client) Client {
	c.client = httpClientFunc
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

func (c client) FromResponse(hr *http.Response) Response {
	response := newResponse(c)
	response.response = hr

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

	return response
}

// Do performs request and returns Response
func (c client) Do(ctx context.Context, target ...interface{}) Response {
	response := newResponse(c)

	request := c.Request().WithContext(ctx)
	request.Close = true

	// call before callback
	if c.before != nil {
		c.before(request)
	}

	// Instantiate http client
	httpClient := c.client()

	var result Response

	// make a http call
	if httpResponse, httpError := httpClient.Do(request); httpError != nil {
		response.error = httpError
		return response
	} else {
		result = c.FromResponse(httpResponse)
	}

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
	strParts := make([]string, 0)

	for _, part := range parts {
		strPart := strings.TrimSpace(fmt.Sprintf("%v", part))
		strPart = strings.TrimLeft(strPart, "/")

		if strPart == "" {
			continue
		}
		strParts = append(strParts, strPart)
	}

	if len(strParts) == 0 {
		return c
	}

	for i, part := range strParts {
		if i < len(strParts)-1 {
			if !strings.HasSuffix(part, "/") {
				part = part + "/"
			}
		}

		c.parts = append(c.parts, part)
	}

	c.base.Path = strings.TrimRight(c.base.Path, "/") + "/"

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
	joined := strings.Join(c.parts, "")

	// Should we append slash?
	if c.appendSlash {
		if joined != "" {
			joined = strings.TrimRight(joined, "/") + "/"
		} else {
			c.base.Path = strings.TrimRight(c.base.Path, "/") + "/"
		}
	}
	p, _ := url.Parse(joined)
	result := c.base.ResolveReference(p)
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
