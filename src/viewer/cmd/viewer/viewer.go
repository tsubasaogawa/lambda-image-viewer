package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"maps"

	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type Image struct {
	Url      string
	Metadata model.Metadata
}

type Headers map[string]string

var (
	//go:embed templates/index.html.tmpl
	tmpl string
)

func main() {
	lambda.Start(Index)
}

func Index(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	p := strings.SplitN(strings.TrimPrefix(r.RawPath, "/"), "/", 2)
	if p == nil || len(p) < 2 {
		msg := "path parsing error. path=" + r.RawPath
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}
	route := p[0]
	key := p[1]

	switch route {
	case "image":
		return generateImageHtml(key)
	case "metadata":
		return generateMetadataJson(key)
	case "cameraroll":
		return generateCamerarollHtml(key)
	default:
		msg := "no route error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}
}

func generateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		msg := "template parsing error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	meta, err := model.New().GetMetadata(getId(key))
	if err != nil {
		log.Println("obtaining metadata error. viewer uses empty metadata.")
		meta = &model.Metadata{}
	}

	image := Image{
		Url:      fmt.Sprintf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key),
		Metadata: *meta,
	}

	w := new(strings.Builder)

	if err = _tmpl.Execute(w, image); err != nil {
		msg := "template execution error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	return responseHtml(w.String(), 200, Headers{}), nil
}

func generateCamerarollHtml(key string) (events.LambdaFunctionURLResponse, error) {
	log.Printf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key)
	return responseHtml("", 200, Headers{"Cache-Control": "private"}), nil
}

func generateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	meta, err := model.New().GetMetadata(getId(key))
	if err != nil {
		msg := fmt.Sprintf(`{"error": "obtaining metadata error, but skips that", "details": "%s"}`, err.Error())
		log.Println(msg)
		return responseJson("{}", 200), nil
		// return responseJson(msg, 500), fmt.Errorf(msg)
	}

	_json, err := json.Marshal(*meta)
	if err != nil {
		msg := `{"error": "json marshal error"}`
		return responseJson(msg, 500), fmt.Errorf(msg)
	}

	return responseJson(string(_json), 200), nil
}

func response(body string, status int, headers Headers) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body:       body,
		Headers:    headers,
		StatusCode: status,
	}
}

func responseHtml(body string, status int, headers Headers) events.LambdaFunctionURLResponse {
	maps.Copy(headers, Headers{"Content-Type": "text/html; charset=utf-8"})
	return response(body, status, headers)
}

func responseJson(body string, status int) events.LambdaFunctionURLResponse {
	return response(body, status, Headers{"Content-Type": "application/json; charset=utf-8"})
}

func getId(key string) string {
	return key
}
