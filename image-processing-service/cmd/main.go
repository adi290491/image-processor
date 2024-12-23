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

	// if err := client.InitS3Service(); err != nil {
	// 	log.Fatal("AWS S3 error:", err)
	// }

	server.Run(":8080")
	// http.HandleFunc("POST /images/transform", handler.HandleTransform)

	// listener, err := net.Listen("tcp", "localhost:8080")

	// if err != nil {
	// 	log.Fatal("Server error:", err)
	// }

	// log.Println("Server running on port:", server.Addr)
	// log.Fatal(server.Serve(listener))

	// InputPath := "assets/input/tiger.jpg"
	// OutputPath := "assets/output/transformed.jpg"

	// transformation := &transformations.Transformation{
	// 	Resize: &transformations.Resize{
	// 		Width:  2000,
	// 		Height: 1800,
	// 	},
	// 	Crop: &transformations.Crop{
	// 		Width:  400,
	// 		Height: 300,
	// 		X:      100,
	// 		Y:      50,
	// 	},
	// 	Rotate: 180,
	// }

	// if err := transformations.Apply(InputPath, OutputPath, *transformation); err != nil {
	// 	log.Fatal("Transformation error:", err)
	// } else {
	// 	log.Println("Transformation successful")
	// }

}
