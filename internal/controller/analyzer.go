package controller

import (
	"strings"

	"golang.org/x/net/html"
)

type Analyzer func(*html.Node, string) (string, error)

func AnalyzerBySelector(selector string) Analyzer {
	if strings.HasPrefix(selector, "xpath:") {
		return HTMLQuery
	}

	return GoQuery
}
