package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net"
	"os"
	"strconv"
)

func getArgs() (int, string) {

	if len(os.Args) != 3 {
		fmt.Printf("Usage: go run client.go <portnumber> <pathfile>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		pathfile := os.Args[2]

		if err != nil {
			fmt.Printf("Usage: go run client.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber, pathfile
		}

	}
	//Should never be reached
	return -1, "issue"
}

func main() {
	port, pathfile := getArgs()
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {
		// image
		defer conn.Close()

		file, _, _ := GetImage(pathfile)
		//encoder := gob.NewEncoder(conn)
		err := png.Encode(conn, file)

		if err != nil {
			fmt.Printf("toto")
			panic(err)
		}

		//decoder := gob.NewDecoder(conn)

		newImage, err := png.Decode(conn)

		if err != nil {
			fmt.Printf("toto1")
			panic(err)
		}

		fg, err := os.Create("./newPhoto.png")
		_ = png.Encode(fg, newImage)

		if err != nil {
			fmt.Printf("toto2")
			panic(err)
		}

		fmt.Printf("success")

		fmt.Printf("#DEBUG MAIN connected\n")

		if err != nil {
			fmt.Printf("toto3")
			fmt.Printf("DEBUG MAIN could not read from server")

			os.Exit(1)
		}

		/*resultString = strings.TrimSuffix(resultString, "\n")
		fmt.Printf("#DEBUG server replied : |%s|\n", resultString)
		time.Sleep(1000 * time.Millisecond)*/

	}

}

func GetImage(filepath string) (image.Image, image.Point, error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	return i, i.Bounds().Size(), err
}
