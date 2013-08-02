package main

import (
	"flag"
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

var (
	red        = flag.Int("red", 0, "red percentage")
	green      = flag.Int("green", 0, "green percentage")
	blue       = flag.Int("blue", 0, "blue percentage")
	inputfile  = flag.String("input", "", "blue percentage")
	outputfile = flag.String("output", "", "blue percentage")
)

func main() {
	flag.Parse()

	file, err := os.Open(*inputfile)
	checkerror(err)

	tofile, err := os.Create(*outputfile)
	checkerror(err)

	defer file.Close()
	defer tofile.Close()

	m, _, err := image.Decode(file)

	bounds := m.Bounds()

	rm := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			if *red < 0 {
				r = invertcolor(r, a)
			}

			if *green < 0 {
				g = invertcolor(g, a)
			}

			if *blue < 0 {
				b = invertcolor(b, a)
			}

			r = r * uint32(*red)
			g = g * uint32(*green)
			b = b * uint32(*blue)

			rm.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	jpeg.Encode(tofile, rm, &jpeg.Options{jpeg.DefaultQuality})
}
