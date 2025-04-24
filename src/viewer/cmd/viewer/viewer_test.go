package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
)

var (
	originalImageGenerator      ImageGenerator
	originalMetadataGenerator   MetadataGenerator
	originalCamerarollGenerator CamerarollGenerator
)

type MockImageGenerator struct{}

func (g *MockImageGenerator) GenerateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "mocked image",
	}, nil
}

type MockMetadataGenerator struct{}

func (g *MockMetadataGenerator) GenerateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "mocked metadata",
	}, nil
}

type MockCamerarollGenerator struct{}

func (g *MockCamerarollGenerator) GenerateCamerarollHtml(pagingKey dynamo.PagingKey) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "mocked cameraroll",
	}, nil
}

func TestIndex(t *testing.T) {
	// Assign mock objects to global variables
	originalImageGenerator = imageGenerator
	originalMetadataGenerator = metadataGenerator
	originalCamerarollGenerator = camerarollGenerator

	imageGenerator = &MockImageGenerator{}
	metadataGenerator = &MockMetadataGenerator{}
	camerarollGenerator = &MockCamerarollGenerator{}

	defer func() {
		imageGenerator = originalImageGenerator
		metadataGenerator = originalMetadataGenerator
		camerarollGenerator = originalCamerarollGenerator
	}()

	testCases := []struct {
		name             string
		rawPath          string
		expectedStatus   int
		expectedBody     string
		expectedErrorMsg string
	}{
		{
			name:           "image route",
			rawPath:        "/image/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "mocked image",
		},
		{
			name:           "metadata route",
			rawPath:        "/metadata/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "mocked metadata",
		},
		{
			name:           "cameraroll route",
			rawPath:        "/cameraroll/",
			expectedStatus: http.StatusOK,
			expectedBody:   "mocked cameraroll",
		},
		{
			name:           "cameraroll route with key",
			rawPath:        "/cameraroll/testid/timestampid",
			expectedStatus: http.StatusOK,
			expectedBody:   "mocked cameraroll",
		},
		{
			name:             "no route",
			rawPath:          "/unknown/test",
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     "no route error",
			expectedErrorMsg: "no route error",
		},
		{
			name:             "path parsing error",
			rawPath:          "/",
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     "path parsing error. path=/",
			expectedErrorMsg: "path parsing error. path=/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := events.LambdaFunctionURLRequest{
				RawPath: tc.rawPath,
			}

			resp, err := Index(r)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			if tc.expectedBody != "" {
				assert.Contains(t, resp.Body, tc.expectedBody)
			}

			if tc.expectedErrorMsg != "" {
				assert.EqualError(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResponseHtml(t *testing.T) {
	body := "test body"
	status := 200
	headers := Headers{"X-Test": "test"}

	resp := responseHtml(body, status, headers)

	assert.Equal(t, body, resp.Body)
	assert.Equal(t, status, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Headers["Content-Type"])
	assert.Equal(t, "test", resp.Headers["X-Test"])
}

func TestResponseJson(t *testing.T) {
	body := "test body"
	status := 200
	headers := Headers{"X-Test": "test"}

	resp := responseJson(body, status, headers)

	assert.Equal(t, body, resp.Body)
	assert.Equal(t, status, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])
	assert.Equal(t, "test", resp.Headers["X-Test"])
}
