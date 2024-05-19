package main

import (
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
)

func TestFillMetadataByExif(t *testing.T) {
	e := &exif.Exif{}

	metadata, err := FillMetadataByExif(e)
	assert.ErrorContains(t, err, `exif: tag "DateTime" is not present`)

	// returns nil if DateTime is none
	assert.Nil(t, metadata)
}
