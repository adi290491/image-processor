package handler

import (
	"fmt"
	transformations "image-processor/pkg"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleTransform(c *gin.Context) {

	// err := r.ParseMultipartForm(math.MaxInt32)

	// if err != nil {
	// 	log.Fatal("Error:", err)
	// }

	// fileHeader, err := c.BindJSON()
	// if err != nil {
	// 	log.Fatal("Request error:", err)
	// }

	uploadDir := "./assets/temp"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {

		c.JSON(500, gin.H{"error": err})
		return
	}

	// tempFilePath := filepath.Join(uploadDir, fileHeader.Filename)

	// err = c.SaveUploadedFile(fileHeader, tempFilePath)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err})
	// 	return
	// }

	var tr transformations.TransformationRequest

	// jsondata := c.PostForm("transformations")
	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(500, gin.H{"Json Error": err})
		return
	}

	response, err := tr.Apply()

	if err != nil {
		HandleError(c, fmt.Errorf("s3 error. Failed to upload image: %v", err), 500)
		return
	}

	c.JSON(200, response)
}
