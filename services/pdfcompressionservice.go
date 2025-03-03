package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CompressFile(inputFile, compressionType string) (string, error) {
	// Define the output file path
	outputDir := "compressed"
	outputFile := filepath.Join(outputDir, filepath.Base(inputFile))

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare the Ghostscript command based on the compression type
	var gsArgs []string
	switch compressionType {
	case "high":
		gsArgs = []string{"-sDEVICE=pdfwrite", "-dCompatibilityLevel=1.4", "-dPDFSETTINGS=/prepress", "-dNOPAUSE", "-dQUIET", "-dBATCH", "-sOutputFile=" + outputFile, inputFile}
	case "medium":
		gsArgs = []string{"-sDEVICE=pdfwrite", "-dCompatibilityLevel=1.4", "-dPDFSETTINGS=/ebook", "-dNOPAUSE", "-dQUIET", "-dBATCH", "-sOutputFile=" + outputFile, inputFile}
	case "low":
		gsArgs = []string{"-sDEVICE=pdfwrite", "-dCompatibilityLevel=1.4", "-dPDFSETTINGS=/screen", "-dNOPAUSE", "-dQUIET", "-dBATCH", "-sOutputFile=" + outputFile, inputFile}
	default:
		return "", fmt.Errorf("unknown compression type: %s", compressionType)
	}

	// Run the Ghostscript command
	cmd := exec.Command("gs", gsArgs...)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ghostscript failed: %w", err)
	}

	// Return the path to the compressed file
	return outputFile, nil
}
