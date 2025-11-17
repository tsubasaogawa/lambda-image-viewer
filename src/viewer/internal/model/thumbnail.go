package model

import (
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

func (t *Table) ListThumbnails(max int64, scanKey dynamo.PagingKey, isPrivate bool) (*[]Thumbnail, dynamo.PagingKey, error) {
	var thumbs []Thumbnail

	// NOTE: GSI の HashKey が Timestamp のみのため、Query を使ってもソート順が保証されない。
	//       このため、cameraroll.go 側でソートを行う。
	//       Timestamp を GSI の SortKey とし、HashKey を固定値にすることで、
	//       Query だけでソートまで完結させることができるが、今回は既存の GSI を利用する。
	query := t.Scan().Index("Timestamp").SearchLimit(max)
	if scanKey != nil {
		query = query.StartFrom(scanKey)
	}

	// Apply filter based on isPrivate flag
	if isPrivate {
		query = query.Filter("contains($, ?)", "Id", "/private/")
	} else {
		query = query.Filter("not contains($, ?)", "Id", "/private/")
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
