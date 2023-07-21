package solace

import (
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	ApiToken *string
	ApiUrl   *string
}

func NewConfig(apiToken *string, apiUrl *string) (*Config, error) {
	var err error

	err = validateApiToken(apiToken)
	if err != nil {
		return nil, err
	}

	err = validateApiUrl(apiUrl)
	if err != nil {
		return nil, err
	}
	var cleanEnvUrl = strings.TrimSuffix(*apiUrl, "/")

	return &Config{
		ApiToken: apiToken,
		ApiUrl:   &cleanEnvUrl,
	}, err
}

func validateApiUrl(envUrl *string) error {
	if envUrl == nil {
		return fmt.Errorf("the API URL is not defined")
	}

	u, err := url.ParseRequestURI(*envUrl)
	if err != nil {
		return fmt.Errorf("the API URL does not seem to be a properly formatted URL")
	}

	if strings.ToLower(u.Scheme) != "https" {
		return fmt.Errorf("use HTTPS protocol for the API URL")
	}

	return nil
}

func validateApiToken(apiToken *string) error {
	if apiToken == nil {
		return fmt.Errorf("the API Token is not defined; to get a token, visit https://docs.solace.com/Cloud/ght_api_tokens.htm and follow the instructions.")
	}

	return nil
}
