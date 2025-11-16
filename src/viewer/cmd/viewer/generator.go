package main

import (
	"github.com/aws/aws-lambda-go/events"
)

type ImageGenerator interface {
	GenerateImageHtml(key string) (events.LambdaFunctionURLResponse, error)
}

type DefaultImageGenerator struct{}

func (g *DefaultImageGenerator) GenerateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	return generateImageHtml(key)
}

type MetadataGenerator interface {
	GenerateMetadataJson(key string) (events.LambdaFunctionURLResponse, error)
}

type DefaultMetadataGenerator struct{}

func (g *DefaultMetadataGenerator) GenerateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	return generateMetadataJson(key)
}
