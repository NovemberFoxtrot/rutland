package main

import (
	"flag"
	"github.com/nfnt/resize"
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
	width = flag.Int("width", 100, "new width")
	height = flag.Int("height", 100, "new height")
)

func outline(source image.Image) image.Image {
	bounds := source.Bounds()

	target := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		r, g, b, a := source.At(bounds.Min.X+1, y).RGBA()
		target.Set(bounds.Min.X, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

		r, g, b, a = source.At(bounds.Max.X-1, y).RGBA()
		target.Set(bounds.Max.X, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
		r, g, b, a := source.At(x, bounds.Min.Y+1).RGBA()
		target.Set(x, bounds.Min.Y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

		r, g, b, a = source.At(x, bounds.Max.Y-1).RGBA()
		target.Set(x, bounds.Max.Y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
			r, g, b, a := source.At(x, y).RGBA()
			target.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return target
}

func mini(source image.Image) image.Image {
	bounds := source.Bounds()

	target := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	yboundary := float64(bounds.Max.Y / 3.0)
	y3 := float64(bounds.Max.Y / 2.0 * 3.0)
	h := float64(bounds.Max.Y)
	times := 30.0

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			r, g, b, _ := source.At(x, y).RGBA()

			var num float64

			if float64(y) < yboundary {
				num = 9.0 * times / h / h * (float64(y) - yboundary) * (float64(y) - yboundary)
			} else if float64(y) < y3 {
				num = 0.0
			} else {
				num = 9.0 * times / h / h * (float64(y) - y3) * (float64(y) - y3)
			}

			target.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(num)})
		}
	}

	return target
}

func smooth(m image.Image) image.Image {
	bounds := m.Bounds()

	target := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y + 1; y <= bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x <= bounds.Max.X-1; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			r1, g1, b1, _ := m.At(x, y-1).RGBA()
			r2, g2, b2, _ := m.At(x-1, y).RGBA()
			r3, g3, b3, _ := m.At(x+1, y).RGBA()
			r4, g4, b4, _ := m.At(x, y+1).RGBA()

			r0 := (r1 + r2 + r3 + r4 + (2.0 * r)) / 6.0
			g0 := (g1 + g2 + g3 + g4 + (2.0 * g)) / 6.0
			b0 := (b1 + b2 + b3 + b4 + (2.0 * b)) / 6.0

			target.Set(x, y, color.RGBA{uint8(r0 / 256), uint8(g0 / 256), uint8(b0 / 256), uint8(a / 256)})
		}
	}

	return target
}

func colour(m image.Image) image.Image {
	bounds := m.Bounds()

	rm := image.NewRGBA(image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			r = invertcolor(r, a)
			g = invertcolor(g, a)
			b = invertcolor(b, a)

			//r = r * uint32(*red)
			//g = g * uint32(*green)
			//b = b * uint32(*blue)

			rm.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return rm
}

func main() {
	flag.Parse()

	file, err := os.Open(*inputfile)
	checkerror(err)

	tofile, err := os.Create(*outputfile)
	checkerror(err)

	defer file.Close()
	defer tofile.Close()

	theImage, _, err := image.Decode(file)

	// theImage = colour(theImage)
	// theImage = smooth(theImage)
	// theImage = outline(theImage)
	// theImage = mini(theImage)

	theImage = resize.Resize(uint(*width), uint(*height), theImage, resize.Lanczos3)

	jpeg.Encode(tofile, theImage, &jpeg.Options{jpeg.DefaultQuality})
}
