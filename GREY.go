package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	pix "pixel"
	"sync"
)

func main() {
	colorTogrey("photo.png", "photogrey.png")
}

func colorTogrey(filepathInput string, filepathOutput string) {
	image := pix.Img{filepathInput}
	tabpix := image.ImageToTab()
	tabgrey := grey(tabpix)
	pix.TabToImage(filepathOutput, tabgrey)

}

func grey(pixels [][]color.Color) (grey [][]color.Color) {

	xLen := len(pixels)
	yLen := len(pixels[0])
	//new image
	newImage := make([][]color.Color, xLen)
	for i := 0; i < len(newImage); i++ {
		newImage[i] = make([]color.Color, yLen)
	}

	//idea is processing pixels in parallel
	wg := sync.WaitGroup{}
	for x := 0; x < xLen; x++ {
		for y := 0; y < yLen; y++ {
			wg.Add(1)
			go func(x, y int) {
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
				wg.Done()
			}(x, y)
		}
	}
	wg.Wait()
	grey = newImage
	return grey
}
