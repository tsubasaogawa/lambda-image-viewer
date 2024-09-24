package model

import (
	"sort"

	"github.com/guregu/dynamo"
)

type Thumbnail struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
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

	sort.Slice(thumbs, func(i, j int) bool { return thumbs[i].Timestamp < thumbs[j].Timestamp })
	return &thumbs, lastKey, nil
}
