package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"watermark/internal"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("ASCII watermark drawer: Please, ensure that your image is located in media/input directory")
	fmt.Println("ASCII watermark drawer: Enter your text: ")

	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	fmt.Println("ASCII watermark drawer: Okay, now enter margin from top to your text: ")
	mt, _ := reader.ReadString('\n')
	marginTop, err := strconv.Atoi(strings.TrimSuffix(mt, "\n"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("ASCII watermark drawer: Okay, now enter margin from left to your text: ")
	ml, _ := reader.ReadString('\n')
	marginLeft, err := strconv.Atoi(strings.TrimSuffix(ml, "\n"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("ASCII watermark drawer: Well, now enter font size to your text (in pixels): ")
	fs, _ := reader.ReadString('\n')
	fontSize, err := strconv.Atoi(strings.TrimSuffix(fs, "\n"))
	if err != nil {
		log.Fatalln(err)
	}

	app, err := internal.NewApp(uint8(marginTop), uint8(marginLeft), uint8(fontSize))
	if err != nil {
		log.Fatalln(err)
	}

	if err = app.WriteText(text); err != nil {
		log.Fatalln(err)
	}
}
