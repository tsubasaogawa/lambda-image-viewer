package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type Image struct {
	Url      string
	Metadata model.Metadata
	Private  bool
	Key      string
}

var (
	//go:embed templates/index.html.tmpl
	imgTmpl string
)

func generateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	_tmpl, err := template.New("index").Parse(imgTmpl)
	if err != nil {
		msg := "template parsing error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	meta, err := getMetadata(key)
	if err != nil {
		log.Println("obtaining metadata error. viewer uses empty metadata.")
		meta = &model.Metadata{}
	}

	imageURL := fmt.Sprintf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key)
	cacheControl := "public, max-age=31536000"
	if model.IsPrivateKey(key) {
		signer, err := newPrivateURLSigner()
		if err != nil {
			log.Printf("Failed to load signing key: %v", err)
			return responseHtml("Internal server error", 500, Headers{"Cache-Control": "private, no-store"}), err
		}
		imageURL, err = signer.Sign(key)
		if err != nil {
			log.Printf("Failed to sign private image: %v", err)
			return responseHtml("Internal server error", 500, Headers{"Cache-Control": "private, no-store"}), err
		}
		cacheControl = "private, no-store"
	}

	image := Image{
		Url:      imageURL,
		Metadata: *meta,
		Private:  model.IsPrivateKey(key),
		Key:      key,
	}

	w := new(strings.Builder)

	if err = _tmpl.Execute(w, image); err != nil {
		msg := "template execution error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": cacheControl}), nil
}
