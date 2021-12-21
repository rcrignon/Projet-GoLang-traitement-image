package main

import (
	img "Img"
	"fmt"
	"image/color"
	"math"
	"runtime"
	"sync"
	"time"
)

var wg3 sync.WaitGroup

func main(){
	start:=time.Now()
	numCPU := runtime.NumCPU()
	test := img.Img{Filepath: "filepath/......../........."}
	m := test.ImageToTab()
	m2 := make([][]color.Color, len(m))
	for i := 0; i < len(m); i++ {
		m2[i] = make([]color.Color, len(m[0]))
	}
	sigma := 10.0
	gauss := matGauss(sigma)
	chunkSize := (len(m) + numCPU - 1) / numCPU
	for i := 0; i < len(m); i += chunkSize {
		end := i + chunkSize
		if end > len(m) {
			end = len(m)
		}
		wg3.Add(1)
		go func(i int, m [][]color.Color) {
			defer wg3.Done()
			for u := i; u < end; u++ {
				for j := 0; j < len(m[0]); j++ {
					mat := matrice7(u, j, m)
					coefpdR := 0.0
					coefpdG := 0.0
					coefpdB := 0.0
					coefg := 0.0
					for k := 0; k < len(mat); k++ {
						for v := 0; v < len(mat[0]); v++ {
							if mat[k][v] != nil {
								r, g, b, _ := mat[k][v].RGBA()
								coef := gauss[k][v]
								coefg = coefg + coef
								coefpdR = coefpdR + coef*(float64(r)/257)
								coefpdG = coefpdG + coef*(float64(g)/257)
								coefpdB = coefpdB + coef*(float64(b)/257)
							} else {
								continue
							}
						}
					}
					m2[u][j] = color.RGBA{R: uint8(coefpdR / coefg), G: uint8(coefpdG / coefg), B: uint8(coefpdB / coefg), A: 255}
				}
			}
		}(i, m)
	}
	wg3.Wait()
	_ = img.TabToImage("filepath/......../.........", m2)
	fmt.Println(time.Since(start))

}

func matrice7(x int, y int, tab [][]color.Color)(m [][]color.Color){
	var m7 = make([][]color.Color, 7)
	for i:=0; i<7; i++{
		var h []color.Color
		for j:=0;j<7;j++{
			if x-3+i>=0 && j+y-3>=0 && i-3+x<len(tab) && j-3+y<len(tab[0]){
				h=append(h,tab[x-3+i][y-3+j])
			}else{
				h=append(h,nil)
			}
		}
		m7[i]=h

	}
	return m7
}

func matGauss(sigma float64)(g [][]float64){
	var gauss = make([][]float64, 7)
	for k:=0; k<7; k++{
		var y []float64
		for v:=0; v<7; v++{
			coef := 1/(2*math.Pi*math.Pow(sigma, 2))*math.Exp(-(math.Pow(math.Abs(float64(3-k)), 2)+(math.Pow(math.Abs(float64(3-v)), 2)))/(2*math.Pow(sigma, 2)))
			y = append(y,coef)
		}
		gauss[k] = y
	}
	return gauss
}