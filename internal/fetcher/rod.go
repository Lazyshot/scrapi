package fetcher

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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
	if err != nil {
		return nil, err
	}

	err = p.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
	})

	return &RodFetcher{
		p: p,
	}, err
}

func (f *RodFactory) Destroy(fetcher Fetcher) {
	// Nothing to do
}
