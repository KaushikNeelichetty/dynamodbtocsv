package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jessevdk/go-flags"
	"os"
)

type options struct {
	TableName string `short:"t" long:"table-name" description:"Name of the dynamo db table" required:"true"`
}

func main() {
	opts := options{}
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(0)
	}

	fmt.Printf("TableName: %v\n", opts.TableName)

	sess, err := session.NewSession()

	if err != nil {
		panic(err)
	}

	dynamoDB := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(opts.TableName),
	}

	result, err := dynamoDB.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(result.Items)

	for result.LastEvaluatedKey != nil {

		params = &dynamodb.ScanInput{
			TableName: aws.String(opts.TableName),
			ExclusiveStartKey: result.LastEvaluatedKey,
		}

		result, err = dynamoDB.Scan(params)

		if err != nil {
			fmt.Println("Query API call failed:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(result.Items)
	}
}
