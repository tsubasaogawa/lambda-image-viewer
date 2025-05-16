package model

import (
	"github.com/guregu/dynamo"
)

type Item struct {
	Id string `json:"id"`
}

func (t *Table) ListItems(max int64, scanKey dynamo.PagingKey) (*[]Item, dynamo.PagingKey, error) {
	var items []Item

	scan := t.Scan().SearchLimit(max)
	if scanKey != nil {
		scan = scan.StartFrom(scanKey)
	}

	lastKey, err := scan.AllWithLastEvaluatedKey(&items)
	if err != nil {
		return nil, nil, err
	}

	return &items, lastKey, nil
}
