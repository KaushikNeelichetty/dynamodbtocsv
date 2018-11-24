package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strings"
)

type options struct {
	TableName string `short:"t" long:"table-name" description:"Name of the dynamo db table" required:"true"`
	Fields    string `short:"f" long:"field-name" description:"List of comma separated fieldName to be output to CSV [Example \"timestamp,count\"] [nested structures will not be flattened]" required:"true"`
}

func main() {
	opts := options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}
	writer := csv.NewWriter(bufio.NewWriter(os.Stdout))
	defer writer.Flush()
	fieldNames := strings.Split(opts.Fields, ",")
	writer.Write(fieldNames)
	dynamoDB := getDynamoDbSession()
	result := doFirstScan(opts, err, dynamoDB, fieldNames, writer)
	doSubSequentScan(result, opts, err, dynamoDB, fieldNames, writer)
}

func doFirstScan(opts options, err error, dynamoDB *dynamodb.DynamoDB, fields []string, writer *csv.Writer) *dynamodb.ScanOutput {
	params := &dynamodb.ScanInput{
		TableName: aws.String(opts.TableName),
	}
	result, err := dynamoDB.Scan(params)
	checkError("Query API call failed", err)
	items := result.Items
	printValues(items, fields, writer)
	return result
}

func doSubSequentScan(result *dynamodb.ScanOutput, opts options, err error, dynamoDB *dynamodb.DynamoDB, fieldsWithType []string, writer *csv.Writer) {
	for result.LastEvaluatedKey != nil {
		params := &dynamodb.ScanInput{
			TableName:         aws.String(opts.TableName),
			ExclusiveStartKey: result.LastEvaluatedKey,
		}
		result, err = dynamoDB.Scan(params)
		checkError("Query API call failed", err)
		items := result.Items
		printValues(items, fieldsWithType, writer)
	}
}

func getDynamoDbSession() *dynamodb.DynamoDB {
	sess, err := session.NewSession()
	checkError("Unable to get dynamodb session", err)
	dynamoDB := dynamodb.New(sess)
	return dynamoDB
}

func getActualValue(iVal interface{}) string {
	if val := fmt.Sprintf("%v", iVal); val != "<nil>" {
		return val
	} else {
		return ""
	}
}

func printValues(items []map[string]*dynamodb.AttributeValue, fields []string, writer *csv.Writer) {
	for _, item := range items {
		valMap := map[string]interface{}{}
		err := dynamodbattribute.UnmarshalMap(item, &valMap)
		checkError("unable to parse dynamodb output", err)
		csvRow := make([]string, len(fields))
		for index, field := range fields {
			csvRow[index] = getActualValue(valMap[field])
		}
		writer.Write(csvRow)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err.Error())
	}
}