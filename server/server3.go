package main

import (
	//"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net"
	"os"
	"strconv"
	"sync"
	//pix"pixel"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup

func getArgs() int {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run server.go <portnumber>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])

		if err != nil {
			fmt.Printf("Usage: go run server.go <portnumber>\n")
			os.Exit(1)

		} else {
			return portNumber
		}

	}
	//PFR should never be reached
	return -1
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
	}

	//If we're here, we did not panic and ln is a valid listener

	connum := 1

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("toto")
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)

		}

		//If we're here, we did not panic and conn is a valid handler to the new connection

		go handleConnection(conn, connum)
		connum += 1

	}
}

func handleConnection(connection net.Conn, connum int) {

	defer connection.Close()
	//decoder := gob.NewDecoder(connection)
	//var Image image.Image
	receiveImage, err := png.Decode(connection)

	if err != nil {
		fmt.Printf("toto1")
		panic(err)
	}

	tab := ImageToTab(receiveImage)
	tabGrey := grey(tab)
	newImage := TabToImage(tabGrey)

	//encoder := gob.NewEncoder(connection)
	err = png.Encode(connection, newImage)
	if err != nil {
		fmt.Printf("toto2")
		panic(err)
	}

	//    if err !=nil{
	//        fmt.Printf("#DEBUG %d handleConnection could not create reader\n", connum)
	//        return
	//    }
	/*for {

		tabpix, err := connReader.ReadString('\n')
		fmt.Printf("toto2")
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}
		tabpix = strings.TrimSuffix(tabpix, "\n")
		greyPix := grey(tabpix)

	}

	/*for {
		inputLine, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}

		//fmt.Printf("#DEBUG RCV |%s|\n", inputLine)
		inputLine = strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)
		splitLine := strings.Split(inputLine, " ")
		returnedString := splitLine[len(splitLine)-1]
		fmt.Printf("#DEBUG %d RCV Returned value |%s|\n", connum, returnedString)
		io.WriteString(connection, fmt.Sprintf("%s\n", returnedString))
	}*/

}

func ImageToTab(im image.Image) [][]color.Color {
	p := im.Bounds().Size()
	var tab = make([][]color.Color, p.X)
	for k := 0; k < p.X; k++ {
		var y []color.Color
		wg.Add(1)
		go func(k int, ligne int, y []color.Color) {
			defer wg.Done()
			for j := 0; j < ligne; j++ {
				y = append(y, im.At(k, j))
			}
			tab[k] = y

		}(k, p.Y, y)

	}
	return tab
}

func TabToImage(tab [][]color.Color) image.Image {

	rect := image.Rect(0, 0, len(tab), len(tab[0]))
	img := image.NewRGBA(rect)
	for x := 0; x < len(tab); x++ {
		wg2.Add(1)
		go func(x int, img *image.RGBA) {
			defer wg2.Done()
			for y := 0; y < len(tab[0]); y++ {
				q := tab[x]
				if q == nil {
					continue
				}
				p := tab[x][y]
				if p == nil {
					continue
				}
				original, ok := color.RGBAModel.Convert(p).(color.RGBA)
				if ok {
					img.Set(x, y, original)
				}
			}
		}(x, img)
	}

	wg2.Wait()

	return img
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
