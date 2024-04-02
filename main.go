package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	inputFolder := "filesInput"   // Path to the input folder containing images
	outputFolder := "filesOutput" // Path to the output folder to save processed images

	// Create output folder if it doesn't exist
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, os.ModePerm)
	}

	// Initialize a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Open the input folder
	err := filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Increment the wait group counter
		wg.Add(1)

		// Process each file in the input folder concurrently
		go func(filePath string) {
			defer wg.Done()
			processFile(filePath, outputFolder)
		}(path)

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func processFile(filePath string, outputFolder string) {
	// Open image file
	img, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer img.Close()

	// Read image data into buffer
	var buf bytes.Buffer
	io.Copy(&buf, img)

	// Create a general holder 256x256
	var holder = [256][256]int{}

	// Iterate through each byte of the image
	for i := 0; i < len(buf.Bytes())-1; i++ {
		holder[buf.Bytes()[i]][buf.Bytes()[i+1]]++
	}

	// Create a grayscale image based on the holder
	newImage := image.NewGray(image.Rect(0, 0, 256, 256))

	// Find max number of repetitions
	max := holder[0][0]
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			if holder[i][j] > max {
				max = holder[i][j]
			}
		}
	}

	// Target average pixel value
	targetAverage := 30

	// Calculate scaling factor to achieve target average
	scaleFactor := float64(targetAverage*256*256) / float64(max)

	// Create an image based on the holder
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			// Scale color intensity based on the count of occurrences
			grayValue := uint8(float64(holder[i][j]) * scaleFactor)
			newImage.SetGray(i, j, color.Gray{Y: grayValue})
		}
	}

	// Generate output file name
	outputFileName := filepath.Join(outputFolder, filepath.Base(filePath)+".jpg")

	// Save the image to output folder
	out, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputFileName, err)
		return
	}
	defer out.Close()

	if err := jpeg.Encode(out, newImage, nil); err != nil {
		fmt.Printf("Error encoding image to file %s: %v\n", outputFileName, err)
	}
}
