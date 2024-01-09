package font2Img

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestToImg(t *testing.T) {
	rgba := GetRGBA("")
	save(rgba)
}

func save(rgba *image.RGBA) {
	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
