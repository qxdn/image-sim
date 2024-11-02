package util

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

func ReadImageFromPath(path string) (*image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadImage(f)
}

func ReadImage(r io.Reader) (*image.Image, error) {
	img, _, err := image.Decode(r)
	return &img, err
}
