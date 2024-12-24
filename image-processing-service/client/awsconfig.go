package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	cfg aws.Config
)

const (
	IMG_BUCKET          = "orig-img-bucket"
	TRANSFORMED_IMG_BKT = "transformation-img-bucket"
	REGION              = "us-east-1"
)

func ConfigureAWS() (err error) {

	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithSharedCredentialsFiles([]string{".env"}),
	)

	if err != nil {
		return fmt.Errorf("AWS configuration error: %v", err)
	}

	_, err = cfg.Credentials.Retrieve(context.TODO())

	if err != nil {
		return fmt.Errorf("AWS error: %v", err)
	}

	return
}

func exitErrorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg+"\n", args...)
}

func uploadImage(reader io.Reader, filename string) (*UploadResponse, error) {

	client := s3.NewFromConfig(cfg)

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(IMG_BUCKET),
		Key:         aws.String(filename),
		Body:        reader,
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(filename))),
	})

	if err != nil {
		err = exitErrorf("File upload error. Failed to upload original image: %v", err)
		return &UploadResponse{nil, ""}, err
	}

	result, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(IMG_BUCKET),
		Key:    aws.String(filename),
	})

	if err != nil {
		err = exitErrorf("Failed to fetch metadata. Error: %v", err)
		return &UploadResponse{nil, ""}, err
	}

	URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", IMG_BUCKET, REGION, filename)
	response := &UploadResponse{
		Metadata: map[string]string{
			"ContentType": aws.ToString(result.ContentType),
			"Size":        fmt.Sprintf("%d", result.ContentLength),
		},
		Url: URL,
	}

	return response, nil
}

func UploadOriginal(imgFile *multipart.FileHeader) (*UploadResponse, error) {

	multipartFile, err := imgFile.Open()
	if err != nil {
		exitErrorf("File IO error: %v", err)
	}

	defer multipartFile.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, multipartFile); err != nil {
		return nil, fmt.Errorf("failed to read multipart file: %w", err)
	}
	return uploadImage(bytes.NewReader(buf.Bytes()), imgFile.Filename)

}

func UploadTransformed(filename string) (*UploadResponse, error) {

	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("File IO error: %v", err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	return uploadImage(file, filepath.Base(fileInfo.Name()))

}

func GetImage(objectKey string) (*s3.GetObjectOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(IMG_BUCKET),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return nil, exitErrorf("Failed to fetch metadata. Error: %v", err)
	}

	return result, nil
}

func ListAllImages(ctx context.Context, pageNum, limit int32) (imageList *ImageListResponse, err error) {
	client := s3.NewFromConfig(cfg)

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(IMG_BUCKET),
	}

	paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = limit
	})

	images := []UploadResponse{}

	currPage := 0
	imageList = &ImageListResponse{}
	for paginator.HasMorePages() && currPage <= int(pageNum) {
		res, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		currPage++
		if currPage != int(pageNum) {
			continue
		}
		imageList.Page = currPage

		imageList.NextToken = *res.NextContinuationToken

		for _, img := range res.Contents {
			images = append(images, UploadResponse{
				Url: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", IMG_BUCKET, REGION, *img.Key),
				Metadata: map[string]string{
					"Size":         fmt.Sprintf("%d", img.Size),
					"LastModified": img.LastModified.String(),
				},
			})
		}

	}

	imageList.Images = images
	return

}
