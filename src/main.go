package main

import (
	"fmt"
	"net/http"
	linkconstants "web-api-searching-and-pagination/src/link-constants"
	"web-api-searching-and-pagination/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Printf("debug: %v\nversion: %v\n", linkconstants.IsDebug(), linkconstants.GetVersion())

	if !linkconstants.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.RedirectTrailingSlash = true

	apiPathPrefix := ""
	engine.GET(apiPathPrefix, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"endpoints": gin.H{
				"dummies": fmt.Sprintf("http://localhost:4000/%sdummies", apiPathPrefix),
			},
		})
	})

	dummyGroup := engine.Group(apiPathPrefix + "dummies")
	routes.RegisterDummyRoutes(dummyGroup)

	port := "4000"
	engine.Run(fmt.Sprintf(":%s", port))

}
