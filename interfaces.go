package rip

import (
	"context"
	"net/http"
	"net/url"
)

// Client
// main object in go-rip. most of methods are chainable, so it provides nice fluent api.
// When Do method is called, Response is returned. This Response is wrapper around http.Response with useful methods.
type Client interface {

	// AppendSlash adds slashes on the end when not presented
	AppendSlash(bool) Client

	// Base sets base url, all rest endpoints are added to this url.
	Base(string) Client

	// Client override default http client.
	Client(func() *http.Client) Client

	// Data sets data to be marshalled to json
	Data(interface{}) Client

	// FromResponse constructs Response from *http.Response
	FromResponse(hr *http.Response) Response

	// Do performs request and returns Response
	Do(ctx context.Context, target ...interface{}) Response

	// Header sets header for request
	Header(key, value string) Client

	// Method sets http method
	Method(method string, parts ...interface{}) Client

	// Path sets path
	Path(parts ...interface{}) Client

	// QueryValues
	QueryValues(values url.Values) Client

	// Request returns prepared request
	Request() *http.Request

	// URL returns actual rest endpoint url
	URL() *url.URL

	// UserAgent overrides default user agent
	UserAgent(agent string) Client

	// Delete HTTP method
	Delete(parts ...interface{}) Client

	// Get HTTP method
	Get(parts ...interface{}) Client

	// Head HTTP method
	Head(parts ...interface{}) Client

	// Options HTTP method
	Options(parts ...interface{}) Client

	// Patch HTTP method
	Patch(parts ...interface{}) Client

	// Put HTTP method
	Put(parts ...interface{}) Client

	// Post HTTP method
	Post(parts ...interface{}) Client
}

// Response is rest response representation. It wraps http.Response, and has body already read.
type Response interface {

	// Body returns response body, if there was an error blank []byte will be returned
	Body() []byte

	// Client returns original client who made the request. This is useful in retry scenarios.
	Client() Client

	// Do calls Do method on client.
	Do(ctx context.Context, target ...interface{}) Response

	// Error returns error if occurred
	Error() error

	// Header returns http header
	Header() http.Header

	// Status fills status into given pointer
	Status(*int) Response

	// Unmarshal unmarshals json response into target
	Unmarshal(target interface{}) Response
}
