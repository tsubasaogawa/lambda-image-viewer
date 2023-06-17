package main

import (
	_ "embed"

	"text/template"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
)

type Photo struct {
	Url string
}

var (
	//go:embed templates/index.html.tmpl
	tmpl string
)

func main() {
	lambda.Start(index)
}

func index(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	param := strings.TrimPrefix(r.Path, "/")
	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	w := new(strings.Builder)
	photo := Photo{Url: "https://photo.ogatube.com/" + param}
	if err = _tmpl.Execute(w, photo); err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
        Body: w.String(),
        Headers: map[string]string{
            "Content-Type": "text/html; charset=utf-8",
        },
        StatusCode: 200,
    }, nil
}
