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

type Env struct {
	EngineHostname string
	EngineSecret   string
}

func NewEnv() (Env, error) {
	engineHostname, err := validateEnv(engineHostnameKey)

	if err != nil {
		return Env{}, err
	}

	engineSecret, err := validateEnv(engineSecretKey)

	if err != nil {
		return Env{}, err
	}

	return Env{EngineHostname: engineHostname, EngineSecret: engineSecret}, nil
}

const engineHostnameKey = "ENGINE_HOSTNAME"
const engineSecretKey = "ENGINE_SECRET"

func handler(event TickerEvent) error {
	env, err := NewEnv()

	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/admin/periodic-work", env.EngineHostname)
	data := PeriodicWorkData(event)
	dataJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", env.EngineSecret))

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

func validateEnv(key string) (string, error) {
	value, present := os.LookupEnv(key)

	if !present {
		return "", fmt.Errorf("%v env var is not present", key)
	}

	return value, nil
}

func main() {
	lambda.Start(handler)
}
