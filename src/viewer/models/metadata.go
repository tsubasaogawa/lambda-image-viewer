package models

import (
	"fmt"
)

type Metadata struct {
	Id          string
	Timestamp   int
	Title       string
	Camera      string
	Lens        string
	Exposure    string
	F           float64
	FocalLength int
}

const (
	REGION = "ap-northeast-1"
	TABLE  = "photo.ogatube.com-item"
)

func (t *Table) GetMetadata(id string) (*Metadata, error) {
	meta := Metadata{}

	if err := t.Get("Id", id).One(&meta); err != nil {
		return nil, fmt.Errorf("foo")
	}

	return &meta, nil
}
