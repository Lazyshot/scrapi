package controller

import (
	"context"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lazyshot/scrapi/internal/fetcher"
	"golang.org/x/net/html"
)

type ScrapeConfig struct {
	Method             string
	URL                string
	DataSelectors      map[string]string
	Multiple           bool
	ItemParentSelector string

	NextPageSelector string
	Limit            int

	VisitItemDetailPage bool
	ItemLinkSelector    string
}

type ScrapeResults struct {
	Num    int
	Data   []map[string]string
	Errors []error
}

type C interface {
	Scrape(ctx context.Context, cfg ScrapeConfig) (ScrapeResults, error)
	AvailableMethods() []fetcher.FetchMethod
}

type Controller struct {
	p *fetcher.FetcherPool
}

func New(p *fetcher.FetcherPool) *Controller {
	return &Controller{p: p}
}

func (c *Controller) Scrape(ctx context.Context, cfg ScrapeConfig) (ScrapeResults, error) {
	res := ScrapeResults{
		Data:   make([]map[string]string, 0),
		Errors: make([]error, 0),
	}
	fetcher, err := c.p.Get(ctx, fetcher.FetchMethod(cfg.Method))
	if err != nil {
		return res, nil
	}
	defer fetcher.Release()

	pageHTML, err := fetcher.Value().Fetch(cfg.URL)
	if err != nil {
		return res, nil
	}

	log.Printf("found html: %s", pageHTML)

	pageDoc, err := html.Parse(strings.NewReader(pageHTML))
	if err != nil {
		return res, err
	}

	if cfg.Multiple {
		for {
			doc := goquery.NewDocumentFromNode(pageDoc)

			if cfg.VisitItemDetailPage {

				childURLs := doc.Find(cfg.ItemParentSelector + " " + cfg.ItemLinkSelector).Map(func(i int, s *goquery.Selection) string {
					v, _ := s.Attr("href")
					return v
				})

				for _, childURL := range childURLs {
					if childURL == "" {
						continue
					}

					childHTML, err := fetcher.Value().Fetch(childURL)
					if err != nil {
						res.Errors = append(res.Errors, err)
						continue
					}

					childDoc, err := html.Parse(strings.NewReader(childHTML))
					if err != nil {
						res.Errors = append(res.Errors, err)
						continue
					}

					if cfg.Limit > 0 && len(res.Data) > cfg.Limit {
						break
					}

					data, err := c.fillData(childDoc, cfg.DataSelectors)
					if err != nil {
						res.Errors = append(res.Errors, err)
						continue
					}

					res.Data = append(res.Data, data)
				}
			} else if cfg.ItemParentSelector != "" {
				parents := doc.Find(cfg.ItemParentSelector)
				log.Printf("found %d parents with %s", parents.Length(), cfg.ItemParentSelector)

				parents.Each(func(i int, s *goquery.Selection) {
					if cfg.Limit > 0 && len(res.Data) > cfg.Limit {
						return
					}

					data, err := c.fillData(s.Get(0), cfg.DataSelectors)
					if err != nil {
						res.Errors = append(res.Errors, err)
						return
					}

					res.Data = append(res.Data, data)
				})
			} else {
				// zip up values into docs?
			}

			if nextPage, ok := doc.Find(cfg.NextPageSelector).Attr("href"); ok {
				pageHTML, err = fetcher.Value().Fetch(nextPage)
				if err != nil {
					return res, nil
				}

				pageDoc, err = html.Parse(strings.NewReader(pageHTML))
				if err != nil {
					return res, err
				}
			} else {
				break
			}
		}

	} else {
		data, err := c.fillData(pageDoc, cfg.DataSelectors)
		if err != nil {
			return res, err
		}

		res.Data = []map[string]string{data}
	}

	res.Num = len(res.Data)

	return res, nil
}

func (c *Controller) fillData(parent *html.Node, selectors map[string]string) (map[string]string, error) {
	data := map[string]string{}

	for k, selector := range selectors {
		v, err := AnalyzerBySelector(selector)(parent, selector)
		if err != nil {
			return nil, err
		}

		data[k] = strings.TrimSpace(v)
	}

	return data, nil
}

func (c *Controller) AvailableMethods() []fetcher.FetchMethod {
	return c.p.Methods()
}
