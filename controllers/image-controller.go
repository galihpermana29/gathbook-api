package controllers

import (
	"learn-golang/models"
	"learn-golang/service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadImages(c *gin.Context) {
	var input models.UploadImageInput

	// binding the multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// handle file uploads

	images := service.SaveUploadedImage(input.Images)

	successResponse := models.ResponseSuccess{
		Data:    images,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	path := filepath.Join("uploads", filename)

	c.File(path)
}
