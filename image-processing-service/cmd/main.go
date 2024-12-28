package main

import (
	"image-processor/client"
	"image-processor/routes"
	"log"

	_ "image-processor/docs"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"


)

// @title           Image Processor API
// @version         1.0
// @description     This is a image processor server that applies various transformations.
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	
	server := gin.Default()


	vips.Startup(nil)
	defer vips.Shutdown()

	routes.RegisterEndpoints(server)
	if err := client.ConfigureAWS(); err != nil {
		log.Fatal("AWS error:", err)
	}

	server.Run(":8080")

}
