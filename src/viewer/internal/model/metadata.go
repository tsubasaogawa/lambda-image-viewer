package model

import (
	"fmt"
)

type Metadata struct {
	Id          string  `json:"id"`
	Timestamp   int64   `json:"timestamp"`
	Title       string  `json:"title"`
	Camera      string  `json:"camera"`
	Lens        string  `json:"lens"`
	Exposure    float64 `json:"exposure"`
	F           float64 `json:"f"`
	FocalLength int     `json:"focal_length"`
	ISO         int     `json:"iso"`
	SS          string  `json:"shutter_speed"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
	IsPrivate   string  `json:"is_private" dynamo:"IsPrivate"`
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
