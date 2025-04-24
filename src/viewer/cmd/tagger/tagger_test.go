package main

import (
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
)

func TestFillMetadataByExif(t *testing.T) {
	e := &exif.Exif{}

	metadata, err := FillMetadataByExif(e)
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
}
