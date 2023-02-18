package fetcher

import (
	"context"

	"github.com/chromedp/chromedp"
)

var _ Fetcher = &ChromeFetcher{}

type ChromeFetcher struct {
	ctx context.Context
}

func (c *ChromeFetcher) Fetch(url string) (string, error) {
	html := ""
	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return "", err
	}

	return html, nil
}

var _ FetcherFactory = &ChromeFactory{}

type ChromeFactory struct {
	ctx context.Context
}

func NewChromeFactory() *ChromeFactory {
	ctx, _ := chromedp.NewContext(context.Background())

	return &ChromeFactory{
		ctx: ctx,
	}
}

func (f *ChromeFactory) Create() (Fetcher, error) {
	ctx, _ := chromedp.NewContext(f.ctx)
	return &ChromeFetcher{
		ctx: ctx,
	}, nil
}

func (f *ChromeFactory) Destroy(fetcher Fetcher) {
	// Nothing to do
}
