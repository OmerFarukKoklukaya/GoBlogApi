package utils

import (
	"errors"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
)

func ImgScaleAndSave(file *multipart.FileHeader) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	var inputFile image.Image
	var outputFile *image.RGBA

	if file.Header["Content-Type"][0] == "image/jpeg" || file.Header["Content-Type"][0] == "image/jpg" {
		inputFile, err = jpeg.Decode(reader)
	} else if file.Header["Content-Type"][0] == "image/png" {
		inputFile, err = png.Decode(reader)

	} else if file.Header["Content-Type"][0] == "image/webp" {
		inputFile, err = webp.Decode(reader)
	} else {
		return errors.New("unsupported file type")
	}
	var newWidth, newHeight = inputFile.Bounds().Dx(), inputFile.Bounds().Dy()

	if inputFile.Bounds().Dx() > 800 {
		ratio := 800 / float32(inputFile.Bounds().Dx())
		newWidth, newHeight = int(float32(inputFile.Bounds().Dx())*ratio), int(float32(inputFile.Bounds().Dy())*ratio)
	} else if inputFile.Bounds().Dy() > 800 {
		ratio := 800 / float32(inputFile.Bounds().Dy())
		newWidth, newHeight = int(float32(inputFile.Bounds().Dx())*ratio), int(float32(inputFile.Bounds().Dy())*ratio)
	}
	outputFile = image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.ApproxBiLinear.Scale(outputFile, outputFile.Bounds(), inputFile, inputFile.Bounds(), draw.Over, nil)

	writer, err := os.Create("./images/" + file.Filename + ".jpeg")
	if err != nil {
		return err
	}
	defer writer.Close()
	jpeg.Encode(writer, outputFile, nil)

	return nil
}
