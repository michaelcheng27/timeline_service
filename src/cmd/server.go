package cmd

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const MOMENTS_TABLE_NAME = "Moments"

type Timeline struct {
	PagingToken *string
	Moments     []Moment
}

type TimelineRequest struct {
	PagingToken *string
}

type Moment struct {
	TimeTaken string
	Id        string
	S3Handle  string
	Latitue   string
	Longitute string
	Camera    string
}

func Serve(request TimelineRequest) (Timeline, error) {
	moments, nextToken, err := getMoments(request.PagingToken)
	if err != nil {
		return Timeline{}, err
	}
	return Timeline{
		PagingToken: nextToken,
		Moments:     moments,
	}, err
}

func getMoments(pagingToken *string) ([]Moment, *string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	var lastEvalutedKey map[string]*dynamodb.AttributeValue
	svc := dynamodb.New(sess)
	if pagingToken != nil {
		log.Infof("pagingToken = %v", *pagingToken)
		json.Unmarshal([]byte(*pagingToken), &lastEvalutedKey)
	}

	params := dynamodb.ScanInput{
		TableName:         aws.String(MOMENTS_TABLE_NAME),
		Limit:             aws.Int64(2),
		ExclusiveStartKey: lastEvalutedKey,
	}
	req, resp := svc.ScanRequest(&params)
	err := req.Send()
	if err != nil {
		return []Moment{}, nil, err
	}
	if resp.Count == nil {
		return []Moment{}, nil, fmt.Errorf("count is nil")
	}
	count := int(*resp.Count)

	if count == 0 {
		return []Moment{}, nil, nil
	}
	moments := make([]Moment, count)
	for i := 0; i < count; i++ {
		err = dynamodbattribute.UnmarshalMap(resp.Items[i], &moments[i])
	}
	nextToken := new(string)
	if resp.LastEvaluatedKey != nil {
		buf, _ := json.Marshal(resp.LastEvaluatedKey)
		*nextToken = string(buf)
	}
	return moments, nextToken, err
}
