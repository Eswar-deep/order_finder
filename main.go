package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/google/uuid"
)

var dbClient *dynamodb.Client
var googleAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")

const tableName = "Orders"

type Order struct {
	OrderID       string  `json:"orderID"`
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
	PreferredTime string  `json:"preferredTime"`
	CreatedAt     string  `json:"createdAt"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	switch {
	case strings.HasSuffix(request.Path, "/order"):
		if request.HTTPMethod == "POST" {
			return handlePost(request)
		} else if request.HTTPMethod == "GET" {
			return handleGet(request)
		}
	case strings.HasSuffix(request.Path, "/orders"):
		if request.HTTPMethod == "GET" {
			return handleListOrders(request)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 405,
		Body:       "Method Not Allowed",
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func handlePost(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input Order

	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid JSON input",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	input.OrderID = uuid.New().String()
	input.CreatedAt = "2025-03-21T08:09:45Z"

	coords, err := geocodeAddress(input.Address)
	if err != nil {
		log.Printf("Geocoding failed: %v", err)
		coords = Coordinates{Lat: 0, Lng: 0}
	}

	input.Lat = coords.Lat
	input.Lng = coords.Lng

	err = insertOrder(input)
	if err != nil {
		log.Printf("Error inserting order: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error storing order",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Order successfully created",
		"orderID": input.OrderID,
		"coords": map[string]float64{
			"lat": coords.Lat,
			"lng": coords.Lng,
		},
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func handleGet(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	orderID := request.QueryStringParameters["orderID"]
	if orderID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing orderID",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"orderID": &types.AttributeValueMemberS{Value: orderID},
		},
	})
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error retrieving order",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Order not found",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	order := map[string]interface{}{}
	for k, v := range result.Item {
		switch val := v.(type) {
		case *types.AttributeValueMemberS:
			order[k] = val.Value
		case *types.AttributeValueMemberN:
			order[k] = val.Value
		}
	}

	body, _ := json.Marshal(order)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func handleListOrders(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	out, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Error scanning orders: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error scanning orders",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}, nil
	}

	var orders []map[string]interface{}

	for _, item := range out.Items {
		order := make(map[string]interface{})
		for k, v := range item {
			switch val := v.(type) {
			case *types.AttributeValueMemberS:
				order[k] = val.Value
			case *types.AttributeValueMemberN:
				order[k] = val.Value
			}
		}
		orders = append(orders, order)
	}

	body, _ := json.Marshal(orders)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func insertOrder(order Order) error {
	_, err := dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"orderID":       &types.AttributeValueMemberS{Value: order.OrderID},
			"name":          &types.AttributeValueMemberS{Value: order.Name},
			"phone":         &types.AttributeValueMemberS{Value: order.Phone},
			"address":       &types.AttributeValueMemberS{Value: order.Address},
			"preferredTime": &types.AttributeValueMemberS{Value: order.PreferredTime},
			"createdAt":     &types.AttributeValueMemberS{Value: order.CreatedAt},
			"lat":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.Lat)},
			"lng":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.Lng)},
		},
	})
	return err
}

func main() {
	lambda.Start(handler)
}
