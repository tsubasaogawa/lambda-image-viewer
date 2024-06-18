package model

import (
	"sort"

	"github.com/guregu/dynamo"
)

type Thumbnail struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

func (t *Table) ListThumbnails(max int64, lastKey dynamo.PagingKey) (*[]Thumbnail, dynamo.PagingKey, error) {
	var thumbs []Thumbnail

	scan := t.Scan().Index("Timestamp").SearchLimit(max)
	if lastKey != nil {
		scan = scan.StartFrom(lastKey)
	}

	lastKey, err := scan.AllWithLastEvaluatedKey(&thumbs)
	if err != nil {
		return nil, nil, err
	}

	sort.Slice(thumbs, func(i, j int) bool { return thumbs[i].Timestamp < thumbs[j].Timestamp })
	return &thumbs, lastKey, nil
}
