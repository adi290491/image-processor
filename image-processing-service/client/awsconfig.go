package client

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
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
		config.WithRegion(REGION),
		config.WithSharedConfigProfile("default"))

	if err != nil {
		return fmt.Errorf("AWS configuration error: %v", err)
	}

	_, err = cfg.Credentials.Retrieve(context.TODO())

	if err != nil {
		return fmt.Errorf("AWS error: %v", err)
	}

	return
}

// func InitS3Service() (err error) {
// 	// svc := s3.New(sess)
// 	// s3.New()

// 	// buckets, err := svc.ListBuckets(nil)

// 	// if err != nil {
// 	// 	return fmt.Errorf("bucket error: %v", err)
// 	// }

// 	// for _, bucket := range buckets.Buckets {
// 	// 	log.Println("Name:", aws.StringValue(bucket.Name), "CreatedAt:", aws.TimeValue(bucket.CreationDate))
// 	// }
// 	// return
// }

func exitErrorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg+"\n", args...)
}

func UploadObject(imgFile *multipart.FileHeader) (*UploadResponse, error) {

	file, err := imgFile.Open()
	if err != nil {
		exitErrorf("File IO error: %v", err)
	}

	defer file.Close()

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(IMG_BUCKET),
		Key:         aws.String(imgFile.Filename),
		Body:        file,
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(imgFile.Filename))),
	})

	if err != nil {
		err = exitErrorf("File upload error. Failed to upload original image: %v", err)
		return &UploadResponse{nil, ""}, err
	}

	result, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(IMG_BUCKET),
		Key:    aws.String(imgFile.Filename),
	})

	if err != nil {
		err = exitErrorf("Failed to fetch metadata. Error: %v", err)
		return &UploadResponse{nil, ""}, err
	}

	URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", IMG_BUCKET, REGION, imgFile.Filename)
	response := &UploadResponse{
		Metadata: map[string]string{
			"ContentType": aws.ToString(result.ContentType),
			"Size":        fmt.Sprintf("%d", result.ContentLength),
		},
		Url: URL,
	}
	// uploader := s3manager.NewUploader(sess)

	// resp, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(ORIG_IMG_BKT),
	// 	Key:    aws.String(imgFile.Filename),
	// 	Body:   file,
	// })

	// if err != nil {
	// 	exitErrorf("File upload error. Failed to upload original image: %v", err)
	// }

	// return resp, nil'
	return response, nil
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

	currPage := 1
	imageList = &ImageListResponse{}
	for paginator.HasMorePages() && currPage <= int(pageNum) {
		res, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if currPage != int(pageNum) {
			continue
		}
		imageList.Page = currPage
		currPage++
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
