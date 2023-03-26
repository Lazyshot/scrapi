package fetcher

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

var _ Fetcher = &RodFetcher{}

type RodFetcher struct {
	p *rod.Page
}

func (c *RodFetcher) Fetch(url string) (string, error) {
	err := c.p.Navigate(url)
	if err != nil {
		return "", err
	}

	return c.p.HTML()
}

var _ FetcherFactory = &RodFactory{}

type RodFactory struct {
	b *rod.Browser
}

func NewRodFactory() *RodFactory {
	return &RodFactory{
		b: rod.New().MustConnect(),
	}
}

func (f *RodFactory) Create() (Fetcher, error) {
	p, err := stealth.Page(f.b)
	return &RodFetcher{
		p: p,
	}, err
}

func (f *RodFactory) Destroy(fetcher Fetcher) {
	// Nothing to do
}
