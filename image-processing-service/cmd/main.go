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

	// server.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowMethods:     []string{"POST", "GET"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	routes.RegisterEndpoints(server)
	if err := client.ConfigureAWS(); err != nil {
		log.Fatal("AWS error:", err)
	}

	server.Run(":8080")

}
