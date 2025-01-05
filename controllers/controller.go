package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"testproject/services"
)

func CompressPdf(c *gin.Context) {
	log.Println("CompressPdf endpoint called")
	c.JSON(http.StatusOK, gin.H{
		"message": "CompressPdf",
	})
}

func Ping(c *gin.Context) {
	log.Println("Ping endpoint called")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UploadAndCompressPDF(c *gin.Context) {
	log.Println("UploadAndCompressPDF endpoint called")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to read the file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	log.Println("File successfully read from the request")

	compressionType := c.DefaultQuery("compression", "medium")
	log.Printf("Compression type: %s", compressionType)

	tempDir := "./uploads"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temp directory: " + err.Error(),
		})
		return
	}

	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("tempfile_%d.pdf", time.Now().UnixNano()))
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Failed to save the uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save the uploaded file: " + err.Error(),
		})
		return
	}
	defer outFile.Close()
	log.Println("Temporary file created")

	_, err = outFile.ReadFrom(file)
	if err != nil {
		log.Printf("Failed to save the uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save the uploaded file: " + err.Error(),
		})
		return
	}
	log.Println("File successfully saved to temporary location")

	resultChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		compressedFilePath, err := services.CompressFile(tempFilePath, compressionType)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- compressedFilePath
	}()

	select {
	case compressedFilePath := <-resultChan:
		log.Println("File successfully compressed")
		compressedFile, err := os.ReadFile(compressedFilePath)
		if err != nil {
			log.Printf("Failed to read the compressed file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read the compressed file: " + err.Error(),
			})
			return
		}

		os.Remove(tempFilePath)
		os.Remove(compressedFilePath)
		log.Println("Temporary files cleaned up")

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"compressed.pdf\""))
		c.Header("Content-Type", "application/pdf")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(compressedFile)
		log.Println("Compressed file sent in response")

	case err := <-errorChan:
		log.Printf("Failed to compress the file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to compress the file: " + err.Error(),
		})

		//case <-time.After(30 * time.Second):
		//	log.Println("File compression timed out")
		//	c.JSON(http.StatusRequestTimeout, gin.H{
		//		"error": "File compression timed out",
		//	})
	}
}
