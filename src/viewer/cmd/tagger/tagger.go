package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func main() {
	lambda.Start(Index)
}

func Index(ctx context.Context, event events.S3Event) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)

	for _, r := range event.Records {
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(r.S3.Bucket.Name),
			Key:    aws.String(r.S3.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}
		e, err := exif.Decode(obj.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", e.String())
		meta := FillMetadataByExif(e)
		fmt.Printf("%#v", *meta)
	}
}

func FillMetadataByExif(e *exif.Exif) *models.Metadata {
	var ts int64 = 0
	if dt, err := e.DateTime(); err == nil {
		ts = dt.Unix()
	}
	return &models.Metadata{
		Id:          "",
		Timestamp:   ts,
		Title:       getExifField(e, exif.ImageDescription).(string), // TODO: The field may be always empty when we use Lightroom
		Camera:      getExifField(e, exif.Model).(string),
		Lens:        getExifField(e, exif.LensModel).(string),
		Exposure:    getExifField(e, exif.ExposureBiasValue).(float64),
		F:           getExifField(e, exif.FNumber).(float64),
		FocalLength: int(getExifField(e, exif.FocalLength).(float64)),
		ISO:         int(getExifField(e, exif.ISOSpeedRatings).(int64)),
		SS:          calcShutterSpeed(getExifField(e, exif.ShutterSpeedValue).(float64)),
	}
}

func getExifField(e *exif.Exif, n exif.FieldName) interface{} {
	tag, err := e.Get(n)
	if err != nil {
		return ""
	}
	fmt.Printf("%s's Format: %d\n", n, tag.Format())

	switch tag.Format() {
	case tiff.StringVal:
		if f, err := tag.StringVal(); err != nil {
			return ""
		} else {
			return f
		}
	case tiff.IntVal:
		if f, err := tag.Int64(0); err != nil {
			return 0
		} else {
			return f
		}
	case tiff.RatVal:
		if nm, dn, err := tag.Rat2(0); err != nil {
			return 0
		} else {
			return float64(nm) / float64(dn)
		}
	default:
		return ""
	}
}
