package rip

import (
	"net/http"
	"net/url"
)

// copyHeaders copies http Header to new map
func copyHeaders(headers http.Header) (result http.Header) {
	result = make(http.Header, len(headers))

	for key := range headers {
		value := headers.Get(key)
		result.Set(key, value)
	}

	return
}

// copyValues copies all url values
func copyValues(values url.Values) (result url.Values) {
	result = url.Values{}

	for key, value := range values {
		for _, val := range value {
			result.Add(key, val)
		}
	}

	return
}
