package web

import (
	"net/url"
)

// URL Decode
func URLDecode(s string) (string, error) {
	r, err := url.QueryUnescape(s)
	return r, err
}
