package main

import (
	"image"
	"image/color"
	_ "image/png"
	"os"
	"runtime"
	"sync"

	"github.com/edsrzf/mmap-go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Canvas")

	filename := "filesInput/randomphoto.png"

	var holder = [256][256]int{}
	var averageValue float64 = 30

	image, _ := processFile(30, &holder)
	imageelement := canvas.NewImageFromImage(&image)
	imageelement.FillMode = canvas.ImageFillOriginal
	imageelement.Resize(fyne.NewSize(100, 100))

	myWindow.SetOnDropped(func(position fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			filename = uri.Path()
			// Open image file
			imgfile, _ := os.Open(filename)
			buf, _ := mmap.Map(imgfile, mmap.RDONLY, 0)
			defer imgfile.Close()

			holder = [256][256]int{}

			numCores := runtime.NumCPU() // Get the number of CPU cores
			chunkSize := len(buf) / numCores
			var wg sync.WaitGroup
			wg.Add(numCores)

			var mu sync.Mutex // Declare a mutex

			// Process chunks based on the number of CPU cores
			processChunk := func(start, end int) {
				defer wg.Done()

				// Use a local holder for each goroutine
				localHolder := [256][256]int{}

				for i := start; i < end-1; i++ {
					localHolder[buf[i]][buf[i+1]]++
				}

				// Merge the results back to the main holder
				mu.Lock()
				for x := 0; x < 256; x++ {
					for y := 0; y < 256; y++ {
						holder[x][y] += localHolder[x][y]
					}
				}
				mu.Unlock()
			}

			// Launch goroutines for each chunk
			for i := 0; i < numCores; i++ {
				start := i * chunkSize
				end := start + chunkSize
				if i == numCores-1 {
					end = len(buf) // Ensure the last chunk processes up to the end of the buffer
				}
				go processChunk(start, end)
			}

			// Wait for all goroutines to finish
			wg.Wait()

			image, _ = processFile(averageValue, &holder)
			imageelement.Image = &image
			imageelement.Refresh()
			buf.Unmap()
			runtime.GC()
		}
	})

	// Add slider
	slider := widget.NewSlider(0.001, 2)
	slider.Step = 0.001
	slider.OnChanged = func(value float64) {
		averageValue = value
		image, _ = processFile(averageValue, &holder)
		imageelement.Image = &image
		imageelement.Refresh()
	}

	// Create the content container
	content := container.NewVBox(
		slider,
		imageelement,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400)) // Set a reasonable window size
	myWindow.ShowAndRun()
}

func processFile(targetAverage float64, holder *[256][256]int) (image.Gray, error) {

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
	return *newImage, nil
}
