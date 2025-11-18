package model

import (
	"strconv"
	"strings"

	"github.com/guregu/dynamo"
)

type Thumbnail struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Private   bool   `json:"private"`
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	IsPrivate string `json:"is_private" dynamo:"IsPrivate"`
}

func (t *Table) ListThumbnails(max int64, scanKey dynamo.PagingKey, isPrivate bool) (*[]Thumbnail, dynamo.PagingKey, error) {
	var thumbs []Thumbnail

	query := t.Get("IsPrivate", strconv.FormatBool(isPrivate)).
		Index("IsPrivate-Timestamp-index").
		Order(dynamo.Descending).
		SearchLimit(max)

	if scanKey != nil {
		query = query.StartFrom(scanKey)
	}

	lastKey, err := query.AllWithLastEvaluatedKey(&thumbs)
	if err != nil {
		return nil, nil, err
	}

	for i := range thumbs {
		thumbs[i].Private = strings.Contains(thumbs[i].Id, "/private/")
	}

	return &thumbs, lastKey, nil
}
