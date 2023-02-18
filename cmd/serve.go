/*
Copyright © 2023 Bryan Peterson <lazyshot@gmail.com>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyshot/scrapi/docs"
	"github.com/lazyshot/scrapi/internal/api"
	"github.com/lazyshot/scrapi/internal/controller"
	"github.com/lazyshot/scrapi/internal/fetcher"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve api",
	Run: func(cmd *cobra.Command, args []string) {
		p := fetcher.NewPool()

		// Register chrome-based fetcher
		err := p.Register("chrome", fetcher.NewChromeFactory(), 10)
		if err != nil {
			panic(err)
		}

		// Register fast stdlib http fetcher
		err = p.Register("http", &fetcher.HTTPFactory{}, 10)
		if err != nil {
			panic(err)
		}

		// create API router
		a := api.New(controller.New(p))

		g := gin.New()

		g.GET("/methods", a.HandleListMethods)
		g.POST("/scrape", a.HandleScrape)

		docs.SwaggerInfo.BasePath = "/"
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		err = g.Run(":8000")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
