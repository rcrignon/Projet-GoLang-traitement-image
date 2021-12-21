package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"runtime"
	"sync"
)
var wg sync.WaitGroup
var wg2 sync.WaitGroup

type Img struct {
	Filepath string
}

func GetImage(filepath string) (image.Image, image.Point, error) {
	f, err := os.Open(filepath)
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	return i, i.Bounds().Size(), err
}

func (im Img) ImageToTab() [][]color.Color {
	numCPU := runtime.NumCPU()
	i, p, _ := GetImage(im.Filepath)
	var tab = make([][]color.Color, p.X)
	for t:=0; t<p.X; t++{
		tab[t]=make([]color.Color, p.Y)
	}
	chunkSize := (len(tab) + numCPU - 1) / numCPU
	for v := 0; v < len(tab); v += chunkSize {
		end := v + chunkSize
		if end > len(tab) {
			end = len(tab)
		}
		wg.Add(1)
		go func(v int, tab [][]color.Color) {
			defer wg.Done()
			for k := v; k < end; k++ {
				for j := 0; j < p.Y; j++ {
					tab[k][j] = i.At(k, j)
				}
			}

		}(v, tab)
	}
	wg.Wait()
	return tab
}

func TabToImage(filepath string, tab [][]color.Color) Img{
	im := Img{Filepath: filepath}
	rect := image.Rect(0,0,len(tab),len(tab[0]))
	img := image.NewRGBA(rect)
	for x:=0; x<len(tab);x++{
		wg2.Add(1)
		go func(x int, img *image.RGBA){
			defer wg2.Done()
			for y:=0; y<len(tab[0]);y++ {
				q:=tab[x]
				if q==nil{
					continue
				}
				p := tab[x][y]
				if p==nil{
					continue
				}
				original, ok := color.RGBAModel.Convert(p).(color.RGBA)
				if ok {
					img.Set(x, y, original)
				}
			}
		}(x,img)
	}

	wg2.Wait()
	fmt.Println(filepath)
	fg,err:= os.Create(filepath)
	_ = png.Encode(fg, img)
	if err!=nil{
		fmt.Println("erreur")	}
	return im
}
