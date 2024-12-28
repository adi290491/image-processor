package routes

import (
	"image-processor/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func RegisterEndpoints(r *gin.Engine) {

	imgHandlers := r.Group("/images")

	imgHandlers.POST("/upload", handler.UploadImage)
	imgHandlers.GET("/", handler.GetImage)
	imgHandlers.GET("/all", handler.ListImages)
	imgHandlers.POST("/transform", handler.HandleTransform)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
