package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jessevdk/go-flags"
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

	fieldNames := strings.Split(opts.Fields, ",")
	printHeaders(fieldNames)

	dynamoDB := getDynamoDbSession()

	result := doFirstScan(opts, err, dynamoDB, fieldNames)

	doSubSequentScan(result, opts, err, dynamoDB, fieldNames)
}

func doFirstScan(opts options, err error, dynamoDB *dynamodb.DynamoDB, fields []string) *dynamodb.ScanOutput {
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
	printValues(items, fields)
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

func printHeaders(fields []string) {
	fmt.Println(strings.Join(fields, ","))
}

func getActualValue(iVal interface{}) string {
	if val := fmt.Sprintf("%v", iVal); val != "<nil>" {
		return val
	} else {
		return ""
	}
}

func printValues(items []map[string]*dynamodb.AttributeValue, fields []string) {
	for _, item := range items {
		valMap := map[string]interface{}{}
		err := dynamodbattribute.UnmarshalMap(item, &valMap)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		outputString := make([]string, len(fields))
		for index, field := range fields {
			outputString[index] = getActualValue(valMap[field])
		}
		fmt.Println(strings.Join(outputString, ","))
	}
}
