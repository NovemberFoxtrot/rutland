package main

import (
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func checkerror(err error) {
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: <input> <output>")
	}

	file, err := os.Open(os.Args[1])
	checkerror(err)

	tofile, err := os.Create(os.Args[2])
	checkerror(err)

	defer file.Close()
	defer tofile.Close()

	m, _, err := image.Decode(file)

	bounds := m.Bounds()

	rm := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, a := m.At(x, y).RGBA()

			rm.Set(x, y, color.RGBA{uint8(a - r), 0, 0, uint8(a)})
		}
	}

	jpeg.Encode(tofile, rm, &jpeg.Options{jpeg.DefaultQuality})
}
