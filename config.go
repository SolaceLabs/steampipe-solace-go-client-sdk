package solace

import (
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	ApiToken  *string
	ApiUrl    *string
	RateLimit *int
}

func NewConfig(apiToken *string, apiUrl *string, rateLimit *int) (*Config, error) {
	var err error

	err = validateRateLimit(rateLimit)
	if err != nil {
		return nil, err
	}

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
		ApiToken:  apiToken,
		ApiUrl:    &cleanEnvUrl,
		RateLimit: rateLimit,
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
		return fmt.Errorf("the API Token is not defined; to get a token, visit the API tab in your Profile page in Make")
	}

	return nil
}

func validateRateLimit(rateLimit *int) error {
	if rateLimit != nil && *rateLimit <= 0 {
		return fmt.Errorf("the rate limit should be a positive number")
	}

	return nil
}
