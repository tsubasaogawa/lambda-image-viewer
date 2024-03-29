package models

import (
	"fmt"
)

type Metadata struct {
	Id          string
	Timestamp   int64
	Title       string
	Camera      string
	Lens        string
	Exposure    float64
	F           float64
	FocalLength int
	ISO         int
	SS          string
}

func (t *Table) GetMetadata(id string) (*Metadata, error) {
	meta := Metadata{}

	if err := t.Get("Id", id).One(&meta); err != nil {
		return nil, fmt.Errorf("GetMetadata() fails")
	}

	return &meta, nil
}

func (t *Table) PutMetadata(meta *Metadata) error {
	if err := t.Put(meta).Run(); err != nil {
		return fmt.Errorf("PutMetadata() fails")
	}

	return nil
}
