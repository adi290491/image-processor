package transformations

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

func resize(imgRef *vips.ImageRef, resizeParam Resize) (err error) {

	width := resizeParam.Width
	height := resizeParam.Height

	scaleX := float64(width) / float64(imgRef.Width())
	scaleY := float64(height) / float64(imgRef.Height())

	scale := min(scaleX, scaleY)
	err = imgRef.Resize(scale, vips.KernelLanczos3)

	if err != nil {
		return
	}

	return
}

func crop(imgRef *vips.ImageRef, cropParam Crop) (err error) {
	// vips.Startup(nil)
	// defer vips.Shutdown()

	left, right, width, height := cropParam.X, cropParam.Y, min(cropParam.Width, imgRef.Width()), min(cropParam.Height, imgRef.Height()) // 100 50 = 50; 100 200 = 100

	if imgRef.Width() <= cropParam.Width {
		width = imgRef.Width() - cropParam.X
	}

	if imgRef.Height() <= cropParam.Height {
		height = imgRef.Height() - cropParam.Y
	}

	err = imgRef.Crop(left, right, width, height)
	if err != nil {
		return
	}

	return
}

func rotate(imgRef *vips.ImageRef, angle float64) (err error) {

	if angle < 0 {
		return fmt.Errorf("invalid rotation angle: %.2f", angle)
	}

	switch angle {
	case 0:
		return imgRef.Rotate(vips.Angle0)
	case 90:
		return imgRef.Rotate(vips.Angle90)
	case 180:
		return imgRef.Rotate(vips.Angle180)
	case 270:
		return imgRef.Rotate(vips.Angle270)
	default:
		return fmt.Errorf("unsupported rotation angle: %.2f, only angles of 0, 90, 180, 270 available", angle)
		// return imgRef.Similarity(1, angle, &vips.ColorRGBA{R: 0, G: 0, B: 0, A: 0}, 0.0, 0.0, 0.0, 0.0)
	}
}
