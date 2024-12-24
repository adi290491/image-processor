package handler

import (
	"context"
	"fmt"
	"image-processor/client"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {

	fileHeader, err := c.FormFile("image")

	if err != nil {
		HandleError(c, fmt.Errorf("image upload error: %w", err), 500)
		return
	}

	response, err := client.UploadOriginal(fileHeader)

	if err != nil {
		HandleError(c, fmt.Errorf("s3 error. Failed to upload image: %v", err), 500)
		return
	}

	c.JSON(200, response)

}

func GetImage(c *gin.Context) {

	objKey := c.Query("key")

	if objKey == "" {
		HandleError(c, fmt.Errorf("object key not present. Key=%s", objKey), 400)
		return
	}

	result, err := client.GetImage(objKey)
	if err != nil {
		HandleError(c, fmt.Errorf("image fetch error: %v", err), 500)
		return
	}
	defer result.Body.Close()
	c.DataFromReader(http.StatusOK, *result.ContentLength, *result.ContentType, result.Body, nil)

}

func ListImages(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		HandleError(c, fmt.Errorf("invalid page number: %v", err), 500)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		HandleError(c, fmt.Errorf("invalid limit: %v", err), 500)
		return
	}

	response, err := client.ListAllImages(context.TODO(), int32(page), int32(limit))

	if err != nil {
		HandleError(c, fmt.Errorf("error while fetching: %v", err), 500)
		return
	}

	c.JSON(http.StatusOK, response)
}
