package transformations

import (
	"errors"
	"fmt"
	"strings"

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

func flip(imgRef *vips.ImageRef, direction string) (err error) {

	if direction == "" {
		return nil
	}

	var dir vips.Direction
	switch strings.ToLower(direction) {
	case "horizontal":
		dir = vips.DirectionHorizontal
	case "vertical":
		dir = vips.DirectionVertical
	default:
		return errors.New("invalid direction - only supports 'horizontal' or 'vertical'")
	}

	return imgRef.Flip(dir)
}

func applyFilters(imgRef *vips.ImageRef, filters *Filters) error {

	/*

		Brightness/Saturation/Hue	Modulate	Brightness 1.0, Saturation 1.0, Hue 0.0
		Gamma Correction	Gamma	1.0 (no correction), <1.0 darken mid-tones
		Sharpness	Sharpen	Sigma: 0.5 to 2.0
		Blur	Blur	Sigma: 1.0 to 5.0
		Contrast Adjustment	Contrast	0.0 (low contrast) to 2.0 (high contrast)
		Grayscale	Modulate or ToColourspace	Saturation 0.0 or InterpretationBW
		Format Conversion	Export, WriteToFile	Use desired file extensions or export params
	*/
	effectsHandlers := map[string]func(*vips.ImageRef, *Filters) error{
		"modulation": func(imgRef *vips.ImageRef, f *Filters) error {
			if f.Modulate != nil {
				imgRef.Modulate(f.Modulate.Brightness, f.Modulate.Saturation, f.Modulate.Hue)
			}
			return nil
		},
		"gamma": func(imgRef *vips.ImageRef, f *Filters) error {

			if f.Gamma == 1.0 {
				return nil
			}
			// if f.Gamma > 1.0 {
			// 	return errors.New("provide a valid value between 1.0 and -1.0")
			// }

			return imgRef.Gamma(f.Gamma)
		},
		"sharpness": func(imgRef *vips.ImageRef, f *Filters) error {
			if f.Sharpness == 0.0 {
				return nil
			}
			return imgRef.Sharpen(f.Sharpness, 1.0, 1.0)
		},
		"blur": func(ir *vips.ImageRef, f *Filters) error {
			if f.Blur != 0.0 {
				return ir.GaussianBlur(f.Blur)
			}
			return nil
		},
		"grayscale": func(ir *vips.ImageRef, f *Filters) error {
			if f.Grayscale {
				return imgRef.ToColorSpace(vips.InterpretationBW)
			}
			return nil
		},
	}

	for f, handler := range effectsHandlers {
		if err := handler(imgRef, filters); err != nil {
			return fmt.Errorf("failed to apply filter: %s due to transformation error:%v", f, err)
		}
	}
	return nil
}
