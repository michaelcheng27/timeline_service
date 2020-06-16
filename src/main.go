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

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(event events.Event, ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	timeline, err := cmd.Serve(cmd.TimelineRequest{})
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
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	log.Info("Hello main")
	lambda.Start(Handler)
}
