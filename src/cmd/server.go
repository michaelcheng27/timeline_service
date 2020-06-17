package cmd

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
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
	TimeTaken      string
	Id             string
	S3Bucket       string
	S3Key          string
	S3PreSignedUrl string
	Latitude       float64
	Longitude      float64
	CameraMake     string
	CameraModel    string
	Width          string
	Height         string
	ImageId        string
}

func Serve(request TimelineRequest) (Timeline, error) {
	moments, nextToken, err := getMoments(request.PagingToken)
	if err != nil {
		return Timeline{}, err
	}
	for i := range moments {
		url, err := getS3PreSignedURL(moments[i])
		if err != nil {
			log.Errorf("faield to generate url")
			continue
		}
		moments[i].S3PreSignedUrl = url
	}
	return Timeline{
		PagingToken: nextToken,
		Moments:     moments,
	}, err
}

func getS3PreSignedURL(moment Moment) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(moment.S3Bucket),
		Key:    aws.String(moment.S3Key),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("The URL is", urlStr)
	return urlStr, err
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
		Limit:             aws.Int64(10),
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

func PutMoment(moment Moment) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	item, err := dynamodbattribute.MarshalMap(moment)
	if err != nil {
		log.Errorf("error = %v", err)
	}
	params := dynamodb.PutItemInput{
		TableName: aws.String(MOMENTS_TABLE_NAME),
		Item:      item,
	}
	req, _ := svc.PutItemRequest(&params)
	err = req.Send()
	return err
}
