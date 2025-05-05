package openlibraryapi

import (
	"net/http"

	"github.com/aczietlow/addToBookshelf/net"
)

const baseURL = "https://openlibrary.org/api"

type Client struct {
	httpClient http.Client
}

type Transport struct {
	UserAgent string
	Transport http.RoundTripper
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Transport: &net.OpenLibraryTransport{
				UserAgent: "Add To Bookshelf/0.1 (aczietlow@gmail.com)",
				Transport: http.DefaultTransport,
			},
		},
	}
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.UserAgent != "" {
		r.Header.Set("User-Agent", t.UserAgent)
	}
	return t.Transport.RoundTrip(r)
}
