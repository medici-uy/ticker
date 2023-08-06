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
	env, err := NewEnv()

	if err != nil {
		return err
	}

	data := PeriodicWorkData(event)
	dataJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	response, err := callPeriodicWork(dataJson, env.EngineHostname, env.EngineSecret)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("engine: %d", response.StatusCode)
	}

	return nil
}

type Env struct {
	EngineHostname string
	EngineSecret   string
}

const engineHostnameKey = "ENGINE_HOSTNAME"
const engineSecretKey = "ENGINE_SECRET"

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

func validateEnv(key string) (string, error) {
	value, present := os.LookupEnv(key)

	if !present {
		return "", fmt.Errorf("%v env var is not present", key)
	}

	return value, nil
}

func callPeriodicWork(dataJson []byte, engineHostname string, engineSecret string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/admin/periodic-work", engineHostname)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", engineSecret))

	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}

func main() {
	lambda.Start(handler)
}
