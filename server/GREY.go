package main

import (
	pix "Img"
	"fmt"
	"image/color"
	_ "image/jpeg"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start:=time.Now()
	colorTogrey("filepath/......../.........", "filepath/......../.........")
	fmt.Println(time.Since(start))
}

func colorTogrey(filepathInput string, filepathOutput string) {
	image := pix.Img{filepathInput}
	tabpix := image.ImageToTab()
	tabgrey := grey(tabpix)
	pix.TabToImage(filepathOutput, tabgrey)

}

func grey(pixels [][]color.Color) (grey [][]color.Color) {
	numCPU := runtime.NumCPU()
	xLen := len(pixels)
	yLen := len(pixels[0])
	newImage := make([][]color.Color, xLen)
	for i := 0; i < len(newImage); i++ {
		newImage[i] = make([]color.Color, yLen)
	}
	chunkSize := (xLen + numCPU - 1) / numCPU
	for i := 0; i < xLen; i += chunkSize {
		end := i + chunkSize
		if end > xLen {
			end = xLen
		}
		wg.Add(1)
		go func(i int, pixels [][]color.Color) {
			defer wg.Done()
			for x := i; x < end; x++ {
				for y := 0; y < yLen; y++ {

					pixel := pixels[x][y]
					originalColor, ok := color.RGBAModel.Convert(pixel).(color.RGBA)
					if !ok {
						fmt.Println("type conversion went wrong")
					}
					grey := uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)
					col := color.RGBA{
						grey,
						grey,
						grey,
						originalColor.A,
					}
					newImage[x][y] = col
				}
			}
		}(i, pixels)
	}
	wg.Wait()
	grey = newImage
	return grey
}
