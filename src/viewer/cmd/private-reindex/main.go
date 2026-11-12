package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type dynamoClient interface {
	ScanPages(*dynamodb.ScanInput, func(*dynamodb.ScanOutput, bool) bool) error
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

func main() {
	table := flag.String("table", os.Getenv("TABLE"), "DynamoDB table name")
	apply := flag.Bool("apply", false, "update matching records instead of performing a dry run")
	flag.Parse()
	if *table == "" {
		log.Fatal("table is required")
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	count, err := reindex(dynamodb.New(sess), *table, *apply)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Matched %d private records (apply=%t)", count, *apply)
}

func reindex(db dynamoClient, table string, apply bool) (int, error) {
	count := 0
	var updateErr error
	err := db.ScanPages(&dynamodb.ScanInput{
		TableName:            aws.String(table),
		ProjectionExpression: aws.String("Id, IsPrivate"),
		ConsistentRead:       aws.Bool(true),
	}, func(page *dynamodb.ScanOutput, _ bool) bool {
		for _, item := range page.Items {
			if item["Id"] == nil || item["Id"].S == nil {
				updateErr = fmt.Errorf("scanned item has no string Id")
				return false
			}
			id := aws.StringValue(item["Id"].S)
			if !model.IsPrivateKey(id) {
				continue
			}
			current := item["IsPrivate"]
			if current != nil && aws.StringValue(current.S) == "true" {
				continue
			}
			if current != nil && (current.S == nil || aws.StringValue(current.S) != "false") {
				updateErr = fmt.Errorf("unexpected IsPrivate value for item %q", id)
				return false
			}
			count++
			if !apply {
				continue
			}
			condition := "attribute_exists(Id) AND attribute_not_exists(IsPrivate)"
			values := map[string]*dynamodb.AttributeValue{
				":new": {S: aws.String("true")},
			}
			if current != nil {
				condition = "attribute_exists(Id) AND IsPrivate = :old"
				values[":old"] = &dynamodb.AttributeValue{S: aws.String("false")}
			}
			_, updateErr = db.UpdateItem(&dynamodb.UpdateItemInput{
				TableName:                 aws.String(table),
				Key:                       map[string]*dynamodb.AttributeValue{"Id": item["Id"]},
				UpdateExpression:          aws.String("SET IsPrivate = :new"),
				ConditionExpression:       aws.String(condition),
				ExpressionAttributeValues: values,
			})
			if updateErr != nil {
				updateErr = fmt.Errorf("update item %q: %w", id, updateErr)
				return false
			}
		}
		return true
	})
	if err != nil {
		return count, fmt.Errorf("scan table: %w", err)
	}
	return count, updateErr
}
