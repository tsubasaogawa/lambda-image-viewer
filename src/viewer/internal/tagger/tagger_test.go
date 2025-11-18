package tagger

import (
	"bytes"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
)

func TestFillMetadataByExif(t *testing.T) {
	e := &exif.Exif{}

	metadata, err := FillMetadataByExif(e, bytes.NewReader(nil))
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
}
