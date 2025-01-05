package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"testproject/services"
)

func CompressPdf(c *gin.Context) {
	// return a test message
	c.JSON(http.StatusOK, gin.H{
		"message": "CompressPdf",
	})
}

func Ping(c *gin.Context) {
	// return a test message
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UploadAndCompressPDF(c *gin.Context) {
	// return a test message
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Taking compression type from the query parameter default to medium
	compressionType := c.DefaultQuery("compression", "medium")

	go func() {
		tempFilePath := "./uploads/tempfile.pdf"
		outFile, err := os.Create(tempFilePath)
		if err != nil {
			return
		}
		defer outFile.Close()
		_, err = outFile.ReadFrom(file)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return
		}

		// Call service to compress file
		compressedFilePath, err := services.CompressFile(tempFilePath, compressionType)
		if err != nil {
			fmt.Println("Error compressing file:", err)
			return
		}

		// Send the result back to the user
		// You can either create a notification system or use a database for tracking.
		// This could be done via an email or system log.

		// For simplicity, printing the result
		fmt.Println("File compressed successfully: ", compressedFilePath)

		// Clean up the temporary file after processing
		err = os.Remove(tempFilePath)
		if err != nil {
			fmt.Println("Error deleting temporary file:", err)
		}
	}()

	// Acknowledge the request has been received
	c.JSON(http.StatusOK, gin.H{"message": "File is being compressed"})
}
