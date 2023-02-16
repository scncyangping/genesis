package util

import (
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
)

func Clip(in io.Reader, out io.Writer, quality int) (err error) {
	err = errors.New("unknow error")
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var origin image.Image
	var fm string
	origin, fm, err = image.Decode(in)
	if err != nil {
		return err
	}

	var canvas image.Image
	var x0, y0, x1, y1 int
	if origin.Bounds().Max.Y > 400 {
		wi := (origin.Bounds().Max.X * 400) / origin.Bounds().Max.Y
		//先缩略
		canvas = resize.Thumbnail(uint(wi), uint(400), origin, resize.Lanczos3)
		x1 = wi - (wi / 100)
		y1 = 360 // 400 - (400 * 0.1)
	} else {
		x1 = origin.Bounds().Max.X
		y1 = origin.Bounds().Max.Y - (origin.Bounds().Max.Y / 10)
		canvas = origin
	}

	switch fm {
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := canvas.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := canvas.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		img := canvas.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	}
	return nil
}
