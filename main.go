package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jessevdk/go-flags"
	"os"
	"strconv"
	"strings"
)

type options struct {
	TableName string `short:"t" long:"table-name" description:"Name of the dynamo db table" required:"true"`
	FieldsWithTypes string `short:"f" long:"fields-with-type" description:"List of comma separated fieldName.DynamoDBType to be output to CSV [Example \"timestamp.S,count.N\"] [nested structures will not be flattened]" required:"true"`
}

func main() {
	opts := options{}
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(0)
	}

	fieldsWithType := strings.Split(opts.FieldsWithTypes, ",")
	printHeaders(fieldsWithType)

	dynamoDB := getDynamoDbSession()

	result := doFirstScan(opts, err, dynamoDB, fieldsWithType)

	doSubSequentScan(result, opts, err, dynamoDB, fieldsWithType)
}

func doFirstScan(opts options, err error, dynamoDB *dynamodb.DynamoDB, fieldsWithType []string) (*dynamodb.ScanOutput) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(opts.TableName),
	}
	result, err := dynamoDB.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	items := result.Items
	printValues(items, fieldsWithType)
	return result
}

func doSubSequentScan(result *dynamodb.ScanOutput, opts options, err error, dynamoDB *dynamodb.DynamoDB, fieldsWithType []string) {
	for result.LastEvaluatedKey != nil {
		params := &dynamodb.ScanInput{
			TableName:         aws.String(opts.TableName),
			ExclusiveStartKey: result.LastEvaluatedKey,
		}
		result, err = dynamoDB.Scan(params)
		if err != nil {
			fmt.Println("Query API call failed:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		items := result.Items
		printValues(items, fieldsWithType)
	}
}

func getDynamoDbSession() *dynamodb.DynamoDB {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	dynamoDB := dynamodb.New(sess)
	return dynamoDB
}

func printHeaders(fieldsWithType []string) {
	header := ""
	for _, fieldWithType := range fieldsWithType {
		field := strings.Split(fieldWithType, ".")[0]
		header = header + field + ","
	}
	fmt.Println(strings.TrimSuffix(header, ","))
}

func printValues(items []map[string]*dynamodb.AttributeValue, fieldsWithType []string) {
	for _, item := range items {
		outputString := ""
		for _, fieldWithType := range fieldsWithType {
			field := strings.Split(fieldWithType, ".")[0]
			dynamoDbType := strings.Split(fieldWithType, ".")[1]
			if fieldValue, ok := item[field]; ok {
				if dynamoDbType == "S" {
					outputString = outputString + *fieldValue.S + ","
				} else if dynamoDbType == "N" {
					outputString = outputString + *fieldValue.N + ","
				} else if dynamoDbType == "B" {
					outputString = outputString + string(fieldValue.B) + ","
				} else if dynamoDbType == "BOOL" {
					outputString = outputString + strconv.FormatBool(*fieldValue.BOOL) + ","
				}
			}
		}
		fmt.Println(strings.TrimSuffix(outputString, ","))
	}
}
