package tagger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // for image.DecodeConfig
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	slambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func Index(ctx context.Context, event events.S3Event) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)

	// If event is empty, fetch all jpg objects from S3.
	// Otherwise, invoke thumbnail generator for the given event.
	if len(event.Records) == 0 {
		log.Println("event is empty. Fetching all objects from S3.")
		r, err := ListAllImageObjects(svc, os.Getenv("ORIGIN_DOMAIN"), "blog/IMG_") // TODO: prefix should be event input
		if err != nil {
			log.Fatal(err)
		}
		event.Records = r
	} else {
		if err := invokeThumbnailGenerator(sess, event); err != nil {
			log.Println(err)
		}
	}

	for _, r := range event.Records {
		if err := ProcessS3Object(svc, r.S3.Bucket.Name, r.S3.Object.Key); err != nil {
			log.Printf("ERROR: failed to process object %s/%s: %v", r.S3.Bucket.Name, r.S3.Object.Key, err)
		}
	}
}

func ProcessS3Object(svc *s3.S3, bucket, key string) error {
	log.Printf("s3://%s/%s\n", bucket, key)
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("could not get object %s/%s: %w", bucket, key, err)
	}
	defer obj.Body.Close()

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		return err
	}

	exifReader := bytes.NewReader(body)
	e, err := exif.Decode(exifReader)
	if err != nil {
		log.Println(err)
		e = &exif.Exif{}
	}

	configReader := bytes.NewReader(body)
	meta, err := FillMetadataByExif(e, configReader)
	if err != nil {
		return err
	}

	meta.Id = key
	meta.IsPrivate = strconv.FormatBool(strings.Contains(key, "/private/"))
	fmt.Printf("%#v", *meta)

	if err := model.New().PutMetadata(meta); err != nil {
		return err
	}

	return nil
}

func FillMetadataByExif(e *exif.Exif, r io.Reader) (*model.Metadata, error) {
	// Get image dimensions
	config, _, err := image.DecodeConfig(r)
	if err != nil {
		log.Println("Failed to decode image config:", err)
		// Assign default values or handle error as appropriate
		config.Width = 0
		config.Height = 0
	}

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
	if _iso, ok := getExifField(e, exif.ISOSpeedRatings).(int64); ok {
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
		Width:       int32(config.Width),
		Height:      int32(config.Height),
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

// ListAllImageObjects retrieves all jpg objects from the specified S3 bucket and prefix.
// It paginates through the results from S3.
func ListAllImageObjects(svc *s3.S3, bucket, prefix string) ([]events.S3EventRecord, error) {
	var records []events.S3EventRecord
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	for {
		// Get a page of objects
		output, err := svc.ListObjectsV2(input)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		// Process objects in the current page
		for _, item := range output.Contents {
			if !strings.HasSuffix(aws.StringValue(item.Key), ".jpg") {
				continue
			}
			record := events.S3EventRecord{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: bucket,
					},
					Object: events.S3Object{
						Key: aws.StringValue(item.Key),
					},
				},
			}
			records = append(records, record)
		}

		// If the result is not truncated, we're done.
		if !aws.BoolValue(output.IsTruncated) {
			break
		}

		// Set the token for the next page request.
		input.ContinuationToken = output.NextContinuationToken
	}

	return records, nil
}
