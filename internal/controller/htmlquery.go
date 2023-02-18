package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

func HTMLQuery(doc *html.Node, selector string) (string, error) {
	selector = strings.TrimPrefix(selector, "xpath:")
	exp, err := xpath.Compile(selector)
	if err != nil {
		return "", err
	}

	res := exp.Evaluate(htmlquery.CreateXPathNavigator(doc))

	switch v := res.(type) {
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), nil
	case xpath.NodeIterator:
		return v.Current().Value(), nil
	}

	return "", errors.New("unknown response type")
}
