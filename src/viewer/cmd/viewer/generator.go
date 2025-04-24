package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
)

type ImageGenerator interface {
	GenerateImageHtml(key string) (events.LambdaFunctionURLResponse, error)
}

type MetadataGenerator interface {
	GenerateMetadataJson(key string) (events.LambdaFunctionURLResponse, error)
}

type CamerarollGenerator interface {
	GenerateCamerarollHtml(pagingKey dynamo.PagingKey) (events.LambdaFunctionURLResponse, error)
}

type DefaultImageGenerator struct{}

func (g *DefaultImageGenerator) GenerateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	return generateImageHtml(key)
}

type DefaultMetadataGenerator struct{}

func (g *DefaultMetadataGenerator) GenerateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	return generateMetadataJson(key)
}

type DefaultCamerarollGenerator struct{}

func (g *DefaultCamerarollGenerator) GenerateCamerarollHtml(pagingKey dynamo.PagingKey) (events.LambdaFunctionURLResponse, error) {
	return generateCamerarollHtml(pagingKey)
}
