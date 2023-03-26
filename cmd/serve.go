/*
Copyright © 2023 Bryan Peterson <lazyshot@gmail.com>
*/
package cmd

import (
	"log"

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
		cf := fetcher.NewChromeFactory()
		if cf.IsChromeInstalled() {
			err := p.Register("chrome", cf, 10)
			if err != nil {
				panic(err)
			}
		} else {
			log.Println("no version of chromium was found installed")
		}

		// Register rod-chrome-based fetcher
		rf := fetcher.NewRodFactory()
		err := p.Register("rod", rf, 10)
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
