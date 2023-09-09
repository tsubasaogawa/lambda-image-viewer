package models

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Table struct {
	dynamo.Table
}

func New() *Table {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})

	return &Table{db.Table(os.Getenv("TABLE"))}
}
