package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	MAX_SCAN_PER_PAGE = 200
)

func main() {
	lambda.Start(Index)
}

func Index(ctx context.Context, event events.S3Event) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)
	tbl := model.New()

	items, lk, err := tbl.ListItems(MAX_SCAN_PER_PAGE, nil)
	if err != nil {
		log.Fatal(err)
	}
	cleanup(svc, tbl, items)

	for len(lk) > 0 {
		items, lk, err = tbl.ListItems(MAX_SCAN_PER_PAGE, lk)
		if err != nil {
			log.Fatal(err)
		}

		cleanup(svc, tbl, items)
	}
}

func cleanup(svc *s3.S3, tbl *model.Table, items *[]model.Item) {
	bucket := os.Getenv("S3_BUCKET_NAME")

	// TODO: goroutine
	for _, item := range *items {
		found, err := exists(svc, item.Id, bucket)
		if err != nil {
			log.Println(err.Error())
		}
		if found || err != nil {
			continue
		}

		d := tbl.Delete("Id", item.Id)
		if err := d.Run(); err != nil {
			log.Printf("Failed to delete item %s: %s\n", item.Id, err.Error())
			continue
		}

		log.Printf("Deleted item %s\n", item.Id)
	}
}

func exists(svc *s3.S3, id string, bucket string) (bool, error) {
	// TODO: use ListObjectsV2 instead of HeadObject
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
	})
	// Object exists (true)
	if err == nil {
		return true, nil
	}

	// Unexpected error and die
	aerr, ok := err.(awserr.Error)
	if !ok {
		return false, err
	}

	// Unexpected error but continue
	if aerr.Code() != s3.ErrCodeNoSuchKey {
		log.Printf("%s: %s\n", id, aerr.Error())
	}

	// Object does not exist (false)
	return false, nil
}
