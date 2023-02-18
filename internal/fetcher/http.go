package fetcher

import (
	"io"
	"net/http"
	"strings"
)

var _ Fetcher = &HTTPFetcher{}

type HTTPFetcher struct {
	c *http.Client
}

func (h *HTTPFetcher) Fetch(url string) (string, error) {
	resp, err := h.c.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sb := new(strings.Builder)
	_, err = io.Copy(sb, resp.Body)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

var _ FetcherFactory = &HTTPFactory{}

type HTTPFactory struct {
}

func (f *HTTPFactory) Create() (Fetcher, error) {
	return &HTTPFetcher{
		c: &http.Client{},
	}, nil
}

func (f *HTTPFactory) Destroy(fetcher Fetcher) {
	// Nothing to do
}
