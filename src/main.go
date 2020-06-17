package main

import (
	"bytes"
	"context"
	"encoding/json"
	"timeline_service/src/cmd"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

// Response is of type APIGatewayProxyResponse since we"re leveraging the
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer
	log.Infof("event = %v", event)
	var request cmd.TimelineRequest
	err := json.Unmarshal([]byte(event.Body), &request)
	timeline, err := cmd.Serve(request)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	body, err := json.Marshal(timeline)
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Headers": "Content-Type, access-control-allow-origin",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST",
		},
	}

	return resp, nil
}

func main() {
	log.Info("Hello main")
	lambda.Start(Handler)
}
