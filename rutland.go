package main

import (
	_ "flag"
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

func invertcolor(somecolor uint32, somealpha uint32) uint32 {
	return somealpha - somecolor
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: <input> <output>")
	}

	var red bool
	var green bool
	var blue bool

	for _, value := range os.Args {
		switch value {
		case "red":
			red = true
		case "green":
			green = true
		case "blue":
			blue = true
		}
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
			r, g, b, a := m.At(x, y).RGBA()

			if red {
				r = invertcolor(r, a)
			}

			if green {
				g = invertcolor(g, a)
			}

			if blue {
				b = invertcolor(b, a)
			}

			rm.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	jpeg.Encode(tofile, rm, &jpeg.Options{jpeg.DefaultQuality})
}
