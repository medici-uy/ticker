package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type TickerEvent struct {
	Frequency string `json:"frequency"`
}

type PeriodicWorkData struct {
	Frequency string `json:"frequency"`
}

func handler(event TickerEvent) error {
	url := fmt.Sprintf("https://%s/admin/periodic-work", os.Getenv("ENGINE_HOSTNAME"))
	data := PeriodicWorkData(event)
	dataJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("ENGINE_SECRET")))

	if err != nil {
		return err
	}

	response, nil := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("engine: %d", response.StatusCode)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
