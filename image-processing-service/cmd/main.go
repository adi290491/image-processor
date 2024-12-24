package main

import (
	"image-processor/client"
	"image-processor/routes"
	"log"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	server := gin.Default()

	routes.RegisterEndpoints(server)
	if err := client.ConfigureAWS(); err != nil {
		log.Fatal("AWS error:", err)
	}

	server.Run(":8080")

}
