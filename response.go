package rip

import (
	"context"
	"encoding/json"
	"net/http"
)

// newResponse returns new Response
func newResponse(client Client) response {
	result := response{
		body:   []byte{},
		client: client,
	}

	return result
}

// response is rest client Response implementation
// It is very lightweight wrapper that serves couple of methods
type response struct {

	// body of response
	body []byte

	// client instance
	client Client

	// store error
	error error

	// response instance
	response *http.Response
}

// Raw sets raw message to given bytes slice
func (r response) Raw(into *[]byte) Response {
	*into = r.body
	return r
}

// Body returns response body
func (r response) Body() (result []byte) {
	return r.body
}

// Client returns client
func (r response) Client() Client {
	return r.client
}

// Do calls Do method on client(for retries)
func (r response) Do(ctx context.Context, target ...interface{}) Response {
	return r.client.Do(ctx, target...)
}

// Error returns error if occurred
func (r response) Error() error {
	return r.error
}

// Header returns header value, if not presented blank string is returned
func (r response) Header() (result http.Header) {
	result = http.Header{}

	if r.response == nil {
		return
	}

	return r.response.Header
}

// Status fills status into given pointer. If there is no response, it's not doing anything
func (r response) Status(status *int) Response {
	if r.response == nil {
		return r
	}

	*status = r.response.StatusCode

	return r
}

// Unmarshal unmarshals json response body into target
func (r response) Unmarshal(target interface{}) Response {
	if r.error != nil {
		return r
	}

	r.error = json.Unmarshal(r.Body(), target)

	return r
}
