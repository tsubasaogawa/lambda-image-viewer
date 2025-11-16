package model

import (
	"sort"
	"strings"

	"github.com/guregu/dynamo"
)

type Thumbnail struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Private   bool   `json:"private"`
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
}

func (t *Table) ListThumbnails(max int64, scanKey dynamo.PagingKey) (*[]Thumbnail, dynamo.PagingKey, error) {
	var thumbs []Thumbnail

	scan := t.Scan().Index("Timestamp").SearchLimit(max)
	if scanKey != nil {
		scan = scan.StartFrom(scanKey)
	}

	lastKey, err := scan.AllWithLastEvaluatedKey(&thumbs)
	if err != nil {
		return nil, nil, err
	}

	for i := range thumbs {
		thumbs[i].Private = strings.Contains(thumbs[i].Id, "/private/")
	}

	sort.Slice(thumbs, func(i, j int) bool { return thumbs[i].Timestamp > thumbs[j].Timestamp })
	return &thumbs, lastKey, nil
}
