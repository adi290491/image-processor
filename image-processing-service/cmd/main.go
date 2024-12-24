package main

import (
	"image-processor/client"
	"image-processor/routes"
	"log"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env error:", err)
	}
	server := gin.Default()

	routes.RegisterEndpoints(server)
	if err := client.ConfigureAWS(); err != nil {
		log.Fatal("AWS error:", err)
	}

	server.Run(":8080")

}
