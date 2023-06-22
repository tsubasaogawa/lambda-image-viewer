package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type widget struct {
	Id        string `dynamo:"id"`
	Timestamp int    `dynamo:"timestamp"`
	Title     string `dynamo:"title"`
	Camera    string `dynamo:"camera"`
	Lens      string `dynamo:"lens"`
	Exposure  string `dynamo:"exposure"`
	F         int    `dynamo:"f"`
}
