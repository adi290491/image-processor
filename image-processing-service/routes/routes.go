package routes

import (
	"image-processor/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(server *gin.Engine) {

	server.POST("/images/upload", handler.UploadImage)
	server.GET("/images", handler.GetImage)
	server.GET("/images/all", handler.ListImages)
	server.POST("/images/transform", handler.HandleTransform)
}
