package transformations

import (
	"fmt"
	"image-processor/client"
	"io"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"
)

type Resize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Crop struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type Transformation struct {
	Resize  *Resize  `json:"resize,omitempty"`
	Crop    *Crop    `json:"crop,omitempty"`
	Rotate  float64  `json:"rotate,omitempty"`
	Flip    string   `json:"flip,omitempty"`
	Filters *Filters `json:"filters,omitempty"`
}

// @Param        transformationRequest body transformations.TransformationRequest true "Transformation Payload"
type TransformationRequest struct {
	Key            string         `json:"image"`
	OutputFileName string         `json:"output"`
	Transformation Transformation `json:"transformation"`
}

type Filters struct {
	Modulate *struct {
		Brightness float64 `json:"brightness"`
		Saturation float64 `json:"saturation"`
		Hue        float64 `json:"hue"`
	} `json:"modulation,omitempty"`
	Gamma     float64 `json:"gamma,omitempty"`
	Sharpness float64 `json:"sharpness,omitempty"`
	Blur      float64 `json:"blur,omitempty"`
	Grayscale bool    `json:"grayscale,omitempty"`
	Format    string  `json:"format,omitempty"`
}

func (t *TransformationRequest) Apply() (response *client.UploadResponse, err error) {

	filename, err := DownloadImage(t.Key)

	if err != nil {
		return nil, fmt.Errorf("file download error: %v", err)
	}

	outputFile, err := t.Transformation.transform(filename, t.OutputFileName)

	if err != nil {
		return nil, fmt.Errorf("transformation error: %v", err)
	}

	response, err = client.UploadTransformed(outputFile)
	if err != nil {
		return nil, fmt.Errorf("upload error: %v", err)
	}

	err = os.Remove(outputFile)
	if err != nil {
		return nil, fmt.Errorf("file deletion error: %v", err)
	}

	err = os.Remove(filename)
	if err != nil {
		return nil, fmt.Errorf("file deletion error: %v", err)
	}

	return
}

func (t *Transformation) transform(input, output string) (outputFile string, err error) {

	imgRef, err := load(input)

	if err != nil {
		return "", fmt.Errorf("image load error:%v", err)
	}

	//do transforms here
	if t.Resize != nil {
		err = resize(imgRef, *t.Resize)
	}

	if err != nil {
		return "", fmt.Errorf("transformation error: %v", err)
	}

	if t.Crop != nil {
		err = crop(imgRef, *t.Crop)
	}
	if err != nil {
		return "", fmt.Errorf("transformation error: %v", err)
	}

	if t.Rotate != 0 {
		err = rotate(imgRef, t.Rotate)
		if err != nil {
			return "", fmt.Errorf("transformation error: %v", err)
		}
	}

	if err = flip(imgRef, t.Flip); err != nil {
		return "", fmt.Errorf("transformation error: %v", err)
	}

	if err = applyFilters(imgRef, t.Filters); err != nil {
		return "", fmt.Errorf("transformation error: %v", err)
	}

	outputFile, err = save(imgRef, output)
	if err != nil {
		return "", fmt.Errorf("image save error: %v", err)
	}

	return
}

func save(imgRef *vips.ImageRef, output string) (outputPath string, err error) {
	// outputPath = filepath.Join("./assets/temp", "processed-"+filepath.Base(input))
	ep := vips.NewDefaultJPEGExportParams()
	imgBytes, _, err := imgRef.Export(ep)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(output, imgBytes, 0644)

	if err != nil {
		return "", err
	}

	return output, nil
}

func load(input string) (imgRef *vips.ImageRef, err error) {

	imgRef, err = vips.NewImageFromFile(input)

	if err != nil {
		return nil, fmt.Errorf("file read error: %v", err)
	}

	return
}

func (t Transformation) String() string {

	return fmt.Sprintf("%v\n%v\n%v", *t.Resize, *t.Crop, t.Rotate)
}

func DownloadImage(key string) (tempFilePath string, err error) {

	resp, err := client.GetImage(key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tempDir := "./assets/temp"

	tempFilePath = filepath.Join(tempDir, key)

	tmpFile, err := os.Create(tempFilePath)

	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)

	if err != nil {
		return "", err
	}
	return
}
