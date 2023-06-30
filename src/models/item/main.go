package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Table struct {
	dynamo.Table
}

type Item struct {
	Id        string
	Timestamp int
	Title     string
	Camera    string
	Lens      string
	Exposure  string
	F         float64
}

const (
	REGION = "ap-northeast-1"
	TABLE  = "photo.ogatube.com-item"
)

func New() *Table {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{
		Region: aws.String(REGION),
	})

	return &Table{db.Table(TABLE)}
}

func (t *Table) _Get(id string) *Item {
	item := Item{}

	if err := t.Get("Id", id).One(&item); err != nil {
		fmt.Errorf("foo")
	}

	return &item
}
