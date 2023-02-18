package controller

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func GoQuery(doc *html.Node, selector string) (string, error) {
	return goquery.NewDocumentFromNode(doc).Find(selector).Text(), nil
}
