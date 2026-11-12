package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type fakeDynamo struct {
	items   []map[string]*dynamodb.AttributeValue
	updates []*dynamodb.UpdateItemInput
	err     error
}

func (d *fakeDynamo) ScanPages(_ *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput, bool) bool) error {
	fn(&dynamodb.ScanOutput{Items: d.items}, true)
	return nil
}

func (d *fakeDynamo) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	d.updates = append(d.updates, input)
	return &dynamodb.UpdateItemOutput{}, d.err
}

func TestReindexStopsOnUpdateFailure(t *testing.T) {
	db := &fakeDynamo{
		items: []map[string]*dynamodb.AttributeValue{
			{"Id": {S: aws.String("private/a.jpg")}, "IsPrivate": {S: aws.String("false")}},
			{"Id": {S: aws.String("private/b.jpg")}, "IsPrivate": {S: aws.String("false")}},
		},
		err: errors.New("conditional check failed"),
	}
	_, err := reindex(db, "table", true)
	if err == nil || len(db.updates) != 1 {
		t.Fatalf("Expected to stop after update failure, got %d updates and error %v", len(db.updates), err)
	}
}

func TestReindex(t *testing.T) {
	db := &fakeDynamo{items: []map[string]*dynamodb.AttributeValue{
		{"Id": {S: aws.String("private/a.jpg")}, "IsPrivate": {S: aws.String("false")}},
		{"Id": {S: aws.String("blog/private/b.jpg")}},
		{"Id": {S: aws.String("blog/private/c.jpg")}, "IsPrivate": {S: aws.String("true")}},
		{"Id": {S: aws.String("public/d.jpg")}, "IsPrivate": {S: aws.String("false")}},
	}}
	count, err := reindex(db, "table", false)
	if err != nil || count != 2 || len(db.updates) != 0 {
		t.Fatalf("Dry run: count=%d, updates=%d, error=%v", count, len(db.updates), err)
	}
	count, err = reindex(db, "table", true)
	if err != nil || count != 2 || len(db.updates) != 2 {
		t.Fatalf("Apply: count=%d, updates=%d, error=%v", count, len(db.updates), err)
	}
	if aws.StringValue(db.updates[0].ConditionExpression) != "attribute_exists(Id) AND IsPrivate = :old" ||
		aws.StringValue(db.updates[1].ConditionExpression) != "attribute_exists(Id) AND attribute_not_exists(IsPrivate)" {
		t.Fatal("Updates must guard against concurrent changes")
	}
}
