package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/tagger"
)

func main() {
	bucket := flag.String("bucket", os.Getenv("ORIGIN_DOMAIN"), "S3 bucket name")
	prefix := flag.String("prefix", "", "S3 object prefix")
	flag.Parse()

	if *bucket == "" {
		log.Fatal("bucket is required")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)

	log.Println("Fetching all objects from S3...")
	records, err := tagger.ListAllImageObjects(svc, *bucket, *prefix)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d objects. Start processing...", len(records))
	for _, r := range records {
		if err := tagger.ProcessS3Object(svc, r.S3.Bucket.Name, r.S3.Object.Key); err != nil {
			log.Printf("ERROR: failed to process object %s/%s: %v", r.S3.Bucket.Name, r.S3.Object.Key, err)
		}
	}
	log.Println("Completed.")
}