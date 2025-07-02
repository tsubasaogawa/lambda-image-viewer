package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	MAX_THUMBNAIL_PER_PAGE = 200
)

var (
	//go:embed templates/camera_roll.html.tmpl
	crTmpl string
)

type CameraRollData struct {
	Thumbnails          *[]model.Thumbnail
	OriginDomain        string
	ViewerDomain        string
	LastKey             string
	PrevKeys            []string
	NextLink            string
	PrevLink            string
	IsPrivate           bool
	SaltForPrivateImage string
}

func generateCamerarollHtml(currentScanKey dynamo.PagingKey, prevKeys []string, isPrivate bool) (events.LambdaFunctionURLResponse, error) {
	m := model.New()
	thumbs, lk, err := m.ListThumbnails(MAX_THUMBNAIL_PER_PAGE, currentScanKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lastEvaluatedKey = %+v\n", lk)

	// Get metadata
	for i, thumb := range *thumbs {
		meta, err := m.GetMetadata(thumb.Id)
		if err != nil {
			log.Printf("Failed to get metadata for %s: %v", thumb.Id, err)
			continue
		}
		(*thumbs)[i].Width = meta.Width
		(*thumbs)[i].Height = meta.Height
	}

	cr, err := template.New("cameraroll").Parse(crTmpl)
	if err != nil {
		log.Fatal(err)
	}

	w := new(strings.Builder)
	// Generate NextLink
	nextLink := ""
	if len(lk) > 0 {
		newPrevKeys := make([]string, len(prevKeys))
		copy(newPrevKeys, prevKeys)
		currentKeyStr := generateLastEvaluatedKeyQueryString(currentScanKey)
		if currentKeyStr != "" {
			newPrevKeys = append(newPrevKeys, currentKeyStr)
		}
		compressedPrevKeys, err := compressPrevKeys(newPrevKeys)
		if err != nil {
			log.Printf("Failed to compress prevKeys for nextLink: %v", err)
			// Fallback to uncompressed or handle error appropriately
			compressedPrevKeys = strings.Join(newPrevKeys, ",")
		}
		nextLink = fmt.Sprintf("/cameraroll/%s?prevKeys=%s", generateLastEvaluatedKeyQueryString(lk), url.QueryEscape(compressedPrevKeys))
		if isPrivate {
			nextLink += "&private=true"
		}
	}

	// Generate PrevLink
	prevLink := ""
	if len(prevKeys) > 0 {
		prevKeyToNavigate := prevKeys[len(prevKeys)-1]  // Get the last key from the history
		remainingPrevKeys := prevKeys[:len(prevKeys)-1] // Remove the last key from the history

		compressedRemainingPrevKeys, err := compressPrevKeys(remainingPrevKeys)
		if err != nil {
			log.Printf("Failed to compress remainingPrevKeys for prevLink: %v", err)
			// Fallback to uncompressed or handle error appropriately
			compressedRemainingPrevKeys = strings.Join(remainingPrevKeys, ",")
		}
		prevLink = fmt.Sprintf("/cameraroll/%s?prevKeys=%s", prevKeyToNavigate, url.QueryEscape(compressedRemainingPrevKeys))
		if isPrivate {
			prevLink += "&private=true"
		}
	} else if len(prevKeys) == 0 && len(currentScanKey) > 0 {
		prevLink = "/cameraroll/"
		if isPrivate {
			prevLink += "?private=true"
		}
	}

	if err = cr.Execute(w, CameraRollData{
		Thumbnails:          thumbs,
		OriginDomain:        os.Getenv("ORIGIN_DOMAIN"),
		ViewerDomain:        os.Getenv("VIEWER_DOMAIN"),
		LastKey:             generateLastEvaluatedKeyQueryString(lk),
		PrevKeys:            prevKeys, // Pass the prevKeys slice to the template
		NextLink:            nextLink,
		PrevLink:            prevLink,
		IsPrivate:           isPrivate,
		SaltForPrivateImage: os.Getenv("SALT_FOR_PRIVATE_IMAGE"),
	}); err != nil {
		return responseHtml("", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}

func generateLastEvaluatedKeyQueryString(lk dynamo.PagingKey) string {
	if len(lk) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s/%s",
		base64.URLEncoding.EncodeToString([]byte(*lk["Id"].S)),
		base64.URLEncoding.EncodeToString([]byte(*lk["Timestamp"].N)),
	)
}

// compressPrevKeys compresses the slice of previous keys into a base64 encoded string.
func compressPrevKeys(prevKeys []string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write([]byte(strings.Join(prevKeys, ",")))
	if err != nil {
		return "", fmt.Errorf("failed to write to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b.Bytes()), nil
}

// decompressPrevKeys decompresses a base64 encoded string into a slice of previous keys.
func decompressPrevKeys(compressed string) ([]string, error) {
	if compressed == "" {
		return []string{}, nil
	}
	data, err := base64.URLEncoding.DecodeString(compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}

	b := bytes.NewReader(data)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("failed to read from gzip reader: %w", err)
	}
	return strings.Split(string(decompressed), ","), nil
}
