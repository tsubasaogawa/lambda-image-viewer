package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	slambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"

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

	if err := invokeThumbnailGenerator(sess, event); err != nil {
		log.Println(err)
	}

	svc := s3.New(sess)
	for _, r := range event.Records {
		log.Printf("s3://%s/%s\n", r.S3.Bucket.Name, r.S3.Object.Key)
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(r.S3.Bucket.Name),
			Key:    aws.String(r.S3.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}
		e, err := exif.Decode(obj.Body)
		if err != nil {
			log.Println(err)
			e = &exif.Exif{}
		}

		meta, err := FillMetadataByExif(e)
		if err != nil {
			log.Fatal(err)
		}

		meta.Id = r.S3.Object.Key
		fmt.Printf("%#v", *meta)

		if err := model.New().PutMetadata(meta); err != nil {
			log.Fatal(err)
		}
	}
}

func FillMetadataByExif(e *exif.Exif) (*model.Metadata, error) {
	ts, err := getLocalUnixtime(e)
	if err != nil {
		log.Println("Timestamp will be set as Now because getLocalUnixtime() got an error: " + err.Error())
	}
	exposure := 0.0
	if _exposure, ok := getExifField(e, exif.ExposureBiasValue).(float64); ok {
		exposure = _exposure
	}
	f := 0.0
	if _f, ok := getExifField(e, exif.FNumber).(float64); ok {
		f = _f
	}
	fl := 0
	if _fl, ok := getExifField(e, exif.FocalLength).(float64); ok {
		fl = int(_fl)
	}
	iso := 0
	if _iso, ok := getExifField(e, exif.ISOSpeedRatings).(int); ok {
		iso = int(_iso)
	}

	ss := "0"
	if ssv, ok := getExifField(e, exif.ShutterSpeedValue).(float64); ok {
		if ss, err = calcShutterSpeed(ssv); err != nil {
			log.Println("Shutter speed will be set as 0 because calcShutterSpeed() got an error: " + err.Error())
		}
	}

	return &model.Metadata{
		Id:          "",
		Timestamp:   ts,
		Title:       getExifField(e, exif.ImageDescription).(string), // TODO: The field may be always empty when we use Lightroom
		Camera:      getExifField(e, exif.Model).(string),
		Lens:        getExifField(e, exif.LensModel).(string),
		Exposure:    exposure,
		F:           f,
		FocalLength: fl,
		ISO:         iso,
		SS:          ss,
	}, nil
}

func getLocalUnixtime(e *exif.Exif) (int64, error) {
	dt, err := e.DateTime()
	if err != nil {
		return time.Now().Unix(), err
	}
	return dt.Unix(), nil
}

func getExifField(e *exif.Exif, n exif.FieldName) interface{} {
	tag, err := e.Get(n)
	if err != nil {
		return ""
	}

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

func invokeThumbnailGenerator(sess *session.Session, event events.S3Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	svc := slambda.New(sess)
	input := &slambda.InvokeInput{
		FunctionName: aws.String(os.Getenv("THUMBNAIL_FUNCTION_NAME")),
		Payload:      payload,
	}
	r, err := svc.Invoke(input)
	if err != nil {
		return err
	}
	log.Println(r)

	return nil
}
