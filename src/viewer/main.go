package main

import (
	_ "embed"

	"fmt"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	key, ok := r.QueryStringParameters["key"]
	if !ok {
		msg := "parameters are invalid"
		return response(msg, 400), fmt.Errorf(msg)
	}

	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		msg := "template parsing error"
		return response(msg, 500), fmt.Errorf(msg)
	}

	w := new(strings.Builder)
	photo := Photo{Url: "https://photo.ogatube.com/" + key}
	if err = _tmpl.Execute(w, photo); err != nil {
		msg := "template execution error"
		return response(msg, 500), fmt.Errorf(msg)
	}

	return response(w.String(), 200), nil
}

func response(body string, status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body: body,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
		},
		StatusCode: status,
	}
}
