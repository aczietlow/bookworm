package net

import "net/http"

type OpenLibraryTransport struct {
	UserAgent string
	Transport http.RoundTripper
}

func (o *OpenLibraryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if o.UserAgent != "" {
		r.Header.Set("User-Agent", o.UserAgent)
	}
	return o.Transport.RoundTrip(r)
}
