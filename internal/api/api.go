package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyshot/scrapi/internal/controller"
)

type API struct {
	c controller.C
}

func New(c controller.C) *API {
	return &API{
		c: c,
	}
}

type ScrapeRequest struct {
	// Method is the name of the fetcher method found at /methods
	Method string `json:"method" binding:"required"`
	// URL is the starting URL (e.g. search result request, something else)
	URL string `json:"url" binding:"required"`

	// DataSelectors are a mapping of field name to css-style selectors or xpath expressions
	// For CSS Style selectors the value is equivalent to a jQuery $(<selector>).text() call
	// For XPath expressions it will be compiled and evaluated
	// XPath expressions must be prefixed with "xpath:"
	DataSelectors map[string]string `json:"dataSelectors"`

	// Multiple signifies this is an array of items
	Multiple bool `json:"multiple"`

	// ItemParentSelector would be the HTML node containing all the
	// data selectors. Data Selectors become relative to the parent
	ItemParentSelector string `json:"itemParentSelector"`

	// Required if pagination is desired
	NextPageSelector string `json:"nextPageSelector"`
	Limit            int    `json:"limit"`

	// visit item detail page and _then_ use data selectors
	VisitItemDetailPage bool   `json:"visitItemDetailPage"`
	ItemLinkSelector    string `json:"itemLinkSelector"`
}

type ScrapeResponse struct {
	Num    int                 `json:"num"`
	Data   []map[string]string `json:"data"`
	Errors []error             `json:"scrapeErrors"`
}

// HandleScrape handle scraping requests
//
//	@Summary		general purpose scraping endpoint
//	@Param			{object}	body	ScrapeRequest	true	"scrape request"
//	@Description	scrapes a website and returns an array of data found given the input config
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ScrapeResponse
//	@Failure		400
//	@Failure		500
//	@Router			/scrape [post]
func (a *API) HandleScrape(c *gin.Context) {
	var req ScrapeRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	cfg := controller.ScrapeConfig{
		Method:             req.Method,
		URL:                req.URL,
		DataSelectors:      req.DataSelectors,
		Multiple:           req.Multiple,
		ItemParentSelector: req.ItemParentSelector,

		NextPageSelector: req.NextPageSelector,
		Limit:            req.Limit,

		VisitItemDetailPage: req.VisitItemDetailPage,
		ItemLinkSelector:    req.ItemLinkSelector,
	}

	res, err := a.c.Scrape(c, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := ScrapeResponse{
		Num:    res.Num,
		Data:   res.Data,
		Errors: res.Errors,
	}

	c.JSON(http.StatusOK, resp)
}

// HandleListMethods lists available and registered fetch methods
//
//	@Summary	list fetch methods
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	string
//	@Router		/methods [get]
func (a *API) HandleListMethods(c *gin.Context) {
	c.JSON(http.StatusOK, a.c.AvailableMethods())
}
