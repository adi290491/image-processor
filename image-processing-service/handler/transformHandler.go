package handler

import (
	"fmt"
	transformations "image-processor/pkg"
	"os"

	"github.com/gin-gonic/gin"
)

// @Summary      Transform Image
// @Description  Apply transformations like resize, crop, or rotate to an image
// @Tags         Images
// @Accept       json
// @Produce      json
// @Param        body body transformations.TransformationRequest true "Transformation Payload"
// @Success      200 {object} client.UploadResponse "Image Info"
// @Failure      500   {object}  APIError "Error Info"
// @Router       /images/transform [post]
func HandleTransform(c *gin.Context) {

	uploadDir := "./assets/temp"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		HandleError(c, fmt.Errorf("error: %v", err), 500)
		return
	}

	var tr transformations.TransformationRequest

	if err := c.ShouldBindJSON(&tr); err != nil {
		HandleError(c, fmt.Errorf("json Error: %v", err), 500)
		return
	}

	response, err := tr.Apply()

	if err != nil {
		HandleError(c, fmt.Errorf("s3 error. Failed to upload image: %v", err), 500)
		return
	}

	c.JSON(200, response)
}
