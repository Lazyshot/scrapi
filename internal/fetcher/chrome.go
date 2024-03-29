package fetcher

import (
	"context"

	cu "github.com/Davincible/chromedp-undetected"
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
	ctx, _, err := cu.New(cu.NewConfig(
		cu.WithHeadless(),
		cu.WithNoSandbox(false),
	))
	if err != nil {
		panic(err)
	}

	return &ChromeFactory{
		ctx: ctx,
	}
}

func (f *ChromeFactory) Create() (Fetcher, error) {
	return &ChromeFetcher{
		ctx: f.ctx,
	}, nil
}

func (f *ChromeFactory) Destroy(fetcher Fetcher) {
	// Nothing to do
}

func (f *ChromeFactory) IsChromeInstalled() bool {
	return chromedp.FromContext(f.ctx).Allocator != nil
}
