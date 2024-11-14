package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/image/draw"
)

const (
	DEFAULT_THUMBNAIL_SIZE = 133
)

func main() {
	lambda.Start(Index)
}

func Index(ctx context.Context, event events.S3Event) {
	/*
			event = {
			  "Records": [{
			    "S3": {
				  "Bucket": { "Name": <bucket_name> },
				  "Object": { "Key": <key> }
				}
			  }]
		    }
	*/
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dl := s3manager.NewDownloader(sess)
	up := s3manager.NewUploader(sess)

	for _, r := range event.Records {
		fo, err := os.CreateTemp("", "oldimg")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(fo.Name())

		n, err := dl.Download(fo, &s3.GetObjectInput{
			Bucket: aws.String(r.S3.Bucket.Name),
			Key:    aws.String(r.S3.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s (%d byte)\n", r.S3.Object.Key, n)

		i, err := generateThumbnail(fo, 0)
		if err != nil {
			log.Fatal(err)
		}
		b := new(bytes.Buffer)
		if err = jpeg.Encode(b, i, &jpeg.Options{Quality: 70}); err != nil {
			log.Fatal(err)
		}

		_, err = up.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(r.S3.Bucket.Name),
			Key:         aws.String(fmt.Sprintf("thumbnail/%s", r.S3.Object.Key)),
			Body:        b,
			ContentType: aws.String("image/jpeg"),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateThumbnail(f *os.File, size int) (*image.RGBA, error) {
	if size <= 0 {
		size = DEFAULT_THUMBNAIL_SIZE
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

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
	return newImage, nil
}
