package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"bytes"

	"github.com/gin-gonic/gin"
	"testproject/services"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
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


func decodeBase64Key(encodedKey string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Base64 key: %v", err)
	}
	return decodedKey, nil
}


func UploadAndCompressPDF(c *gin.Context) {
	log.Println("UploadAndCompressPDF endpoint called")

	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	// Check if the request is an OPTIONS preflight request
	if c.Request.Method == "OPTIONS" {
		c.JSON(http.StatusOK, gin.H{"message": "Preflight request successful"})
		return
	}

	// Step 1: Read the encrypted file from the request
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

	// Step 2: Decrypt the file (decrypt the file using a Base64 key)
	encodedKey := "UG+zlV6pEODX2mZ8TFMJn5DaWK8SUCjoZl3gSg5G6WE=" // Assuming the key is passed as a query parameter
	if encodedKey == "" {
		log.Println("Missing encryption key")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Encryption key missing in request",
		})
		return
	}

	key, err := decodeBase64Key(encodedKey)
	if err != nil {
		log.Printf("Failed to decode the encryption key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode the encryption key: " + err.Error(),
		})
		return
	}

	// Step 3: Use a buffer to read file data
	var encryptedData bytes.Buffer
	_, err = io.Copy(&encryptedData, file) // Write file data to the buffer
	if err != nil {
		log.Printf("Failed to read file contents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read the file contents: " + err.Error(),
		})
		return
	}

	// Step 4: Decrypt the file data using AES-GCM
	decryptedData, err := decryptAESGCM(encryptedData.Bytes(), key)
	if err != nil {
		log.Printf("Failed to decrypt the file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to decrypt the file: " + err.Error(),
		})
		return
	}

	// Step 5: Save the decrypted file temporarily
	tempDir := "./uploads"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temp directory: " + err.Error(),
		})
		return
	}

	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("decryptedfile_%d.pdf", time.Now().UnixNano()))
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary file: " + err.Error(),
		})
		return
	}
	defer outFile.Close()

	_, err = outFile.Write(decryptedData)
	if err != nil {
		log.Printf("Failed to save decrypted file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save decrypted file: " + err.Error(),
		})
		return
	}
	log.Println("Decrypted file saved")

	// Step 6: Compress the file
	compressionType := c.DefaultQuery("compression", "medium")
	log.Printf("Compression type: %s", compressionType)

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

		// Step 7: Encrypt the compressed file
		compressedFile, err := os.ReadFile(compressedFilePath)
		if err != nil {
			log.Printf("Failed to read the compressed file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read the compressed file: " + err.Error(),
			})
			return
		}

		// Encrypt the compressed file before sending it back
		encryptedCompressedFile, err := encryptAESGCM(compressedFile, key)
		if err != nil {
			log.Printf("Failed to encrypt the compressed file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to encrypt the compressed file: " + err.Error(),
			})
			return
		}

		// Clean up temporary files
		os.Remove(tempFilePath)
		os.Remove(compressedFilePath)
		log.Println("Temporary files cleaned up")

		// Send the encrypted compressed file in the response
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"compressed_encrypted.pdf\""))
		c.Header("Content-Type", "application/octet-stream")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encryptedCompressedFile)
		log.Println("Encrypted compressed file sent in response")

	case err := <-errorChan:
		log.Printf("Failed to compress the file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to compress the file: " + err.Error(),
		})

	case <-time.After(30 * time.Second):
		log.Println("File compression timed out")
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": "File compression timed out",
		})
	}
}

func encryptAESGCM(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decryptAESGCM(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
