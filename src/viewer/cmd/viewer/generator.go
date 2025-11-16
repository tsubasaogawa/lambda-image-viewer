package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
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

type CamerarollGenerator interface {
	CameraRollHandler(currentScanKey dynamo.PagingKey, prevKeys []string, isPrivate bool) (events.LambdaFunctionURLResponse, error)
}

type DefaultCamerarollGenerator struct {
	DB DB
}
