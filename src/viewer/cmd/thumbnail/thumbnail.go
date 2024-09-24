package main

import (
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

const (
	DEFAULT_THUMBNAIL_SIZE = 133
)

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func SquareTrimImage(img image.Image, size int) image.Image {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	shorter := width
	if height < shorter {
		shorter = height
	}

	top := (height - shorter) / 2
	left := (width - shorter) / 2

	newImage := image.NewRGBA(image.Rect(0, 0, size, size))

	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, image.Rect(left, top, width-left, height-top), draw.Over, nil)

	return newImage
}
