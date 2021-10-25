package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

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

func (im Img) ImageToTab() [][]color.Color{
	i, p, _ :=GetImage(im.Filepath)
	var tab [][]color.Color
	for k:=0; k<p.X; k++{
		var y []color.Color
		for j:=0; j<p.Y;j++{
			y = append(y,i.At(k,j))
		}
		tab = append(tab,y)
	}
	return tab
}

func TabToImage(filepath string, tab [][]color.Color) Img{
	im := Img{Filepath: filepath}
	rect := image.Rect(0,0,len(tab),len(tab[0]))
	img := image.NewRGBA(rect)
	for x:=0; x<len(tab);x++{
		for y:=0; y<len(tab[0]);y++ {
			q:=tab[x]
			if q==nil{
				continue
			}
			p := tab[x][y]
			if p==nil{
				continue
			}
			original,ok := color.RGBAModel.Convert(p).(color.RGBA)
			if ok{
				img.Set(x,y,original)
			}
		}
	}
	fmt.Println(filepath)
	fg,err:= os.Create(filepath)
	_ = png.Encode(fg, img)
	if err!=nil{
		fmt.Println("erreur")
	}
	return im
}
