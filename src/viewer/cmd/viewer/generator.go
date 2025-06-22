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
	GenerateCamerarollHtml(currentScanKey dynamo.PagingKey, prevKeys []string, isPrivate bool) (events.LambdaFunctionURLResponse, error)
	DecompressPrevKeys(compressed string) ([]string, error) // Add this line
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

func (g *DefaultCamerarollGenerator) GenerateCamerarollHtml(currentScanKey dynamo.PagingKey, prevKeys []string, isPrivate bool) (events.LambdaFunctionURLResponse, error) {
	return generateCamerarollHtml(currentScanKey, prevKeys, isPrivate)
}

func (g *DefaultCamerarollGenerator) DecompressPrevKeys(compressed string) ([]string, error) {
	return decompressPrevKeys(compressed)
}
