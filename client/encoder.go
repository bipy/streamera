package client

import (
	"bytes"
	"gocv.io/x/gocv"
	"image/jpeg"
)

const jpegQuality = 50

var jpegOption = &jpeg.Options{Quality: jpegQuality}

func encodeImage(mat *gocv.Mat) []byte {
	image, err := mat.ToImage()
	if err != nil {
		return nil
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, jpegOption)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
