package openlibraryapi

import (
	"net/http"
	"time"

	"github.com/aczietlow/bookworm/pkg/bookcache"
)

const baseURL = "https://openlibrary.org"

type Client struct {
	httpClient http.Client
	cache      bookcache.Cache
}

type Transport struct {
	UserAgent string
	Transport http.RoundTripper
}

func NewClient(httpTimeout time.Duration, cacheTTL time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Transport: &Transport{
				UserAgent: "Bookworm/0.1 (aczietlow@gmail.com)",
				Transport: http.DefaultTransport,
			},
			Timeout: httpTimeout,
		},
		cache: bookcache.NewCacheStorage(cacheTTL),
	}
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.UserAgent != "" {
		r.Header.Set("User-Agent", t.UserAgent)
	}
	return t.Transport.RoundTrip(r)
}
