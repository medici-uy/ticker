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

const engineHostnameKey = "ENGINE_HOSTNAME"
const engineSecretKey = "ENGINE_SECRET"

func handler(event TickerEvent) error {
	engineHostname, present := os.LookupEnv(engineHostnameKey)

	if !present {
		return fmt.Errorf("%v env var is not present", engineHostnameKey)
	}

	engineSecret, present := os.LookupEnv(engineHostnameKey)

	if !present {
		return fmt.Errorf("%v env var is not present", engineSecretKey)
	}

	url := fmt.Sprintf("https://%s/admin/periodic-work", engineHostname)
	data := PeriodicWorkData(event)
	dataJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", engineSecret))

	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)

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
