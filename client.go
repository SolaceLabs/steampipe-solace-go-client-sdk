package solacesdk

import (
	"bytes"
	"context"
	"strconv"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	defaultPageSize  = 100
	defaultRateLimit = 50
)

type Client struct {
	client      *http.Client
	rateLimiter <-chan time.Time
	baseUrl     string
	apiToken    string
}

func GetClient(config *Config) *Client {
	if config.RateLimit == nil {
		config.RateLimit = &defaultRateLimit
	}

	// rate limiter with 20% burstable rate
	rateLimiter := make(chan time.Time, *config.RateLimit/20)
	go func() {
		for t := range time.Tick(time.Minute / time.Duration(*config.RateLimit)) {
			rateLimiter <- t
		}
	}()

	clientInstance := &Client{
		client:      http.DefaultClient,
		rateLimiter: rateLimiter,
		apiToken:    *config.ApiToken,
		baseUrl:     *config.ApiUrl,
	}

	return clientInstance
}

func (at *Client) rateLimit() {
	<-at.rateLimiter
}

func (at *Client) Get(config *RequestConfig, target interface{}) ([]byte, error) {
	at.rateLimit()

	// prepare the request URL
	req, err := at.createAuthorizedRequest(fmt.Sprintf("%s/api/v2/%s", at.baseUrl, config.Endpoint), http.MethodGet)
	if err != nil {
		return nil, err
	}
	at.setQueryParams(req, config)

	// do the call
	b, err := at.do(req)
	if err != nil {
		return nil, err
	}

	// parse the body
	if target != nil {
		err = json.Unmarshal(b, target)
		if err != nil {
			return nil, fmt.Errorf("JSON decode failed: %s error: %w", b, err)
		}
	}

	return b, nil
}

func (at *Client) Post(config *RequestConfig, target interface{}) ([]byte, error) {
	at.rateLimit()

	// prepare the request URL
	var body, err = json.Marshal(config.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal body: %w", err)
	}
	req, err := at.createAuthorizedRequest(fmt.Sprintf("%s/api/v2/%s", at.baseUrl, config.Endpoint), http.MethodPost)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	// do the call
	b, err := at.do(req)
	if err != nil {
		return nil, err
	}

	// parse the body
	if target != nil {
		err = json.Unmarshal(b, target)
		if err != nil {
			return nil, fmt.Errorf("JSON decode failed: %s error: %w", b, err)
		}
	}

	return b, nil
}

func (at *Client) createAuthorizedRequest(apiUrl string, method string) (*http.Request, error) {
	log.Println(fmt.Sprintf("[MAKE]:[%s] -> %s", method, apiUrl))

	// make a new request
	req, err := http.NewRequestWithContext(context.Background(), method, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	// set headers and query params
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiToken))

	return req, nil
}

func (at *Client) setQueryParams(req *http.Request, config *RequestConfig) {
	// set pagination params
	if config.Pagination != nil && config.Pagination.paginate {
		config.Params.Set("pageNumber", strconv.Itoa(config.Pagination.pageNumber))
		config.Params.Set("pageSize", strconv.Itoa(config.Pagination.pageSize))
	}

	// encode params
	req.URL.RawQuery = config.Params.Encode()
	log.Println(fmt.Sprintf("[SOLACE]:[PARAMS] %s", req.URL.RawQuery))
}

func (at *Client) do(req *http.Request) ([]byte, error) {
	var reqUrl = req.URL.RequestURI()

	// do the call
	var resp, err = at.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failure on %s: %w", reqUrl, err)
	}
	defer resp.Body.Close()

	// handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, NewHttpError(reqUrl, resp)
	}

	// read response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP Read error on response for %s: %w", reqUrl, err)
	}

	return b, nil
}

func (at *Client) handleKnownErrors(err error) error {
	var httpErr = getHttpError(err)
	if httpErr == nil {
		return err
	}

	// 403 Forbidden or 404 Not Found
	switch httpErr.StatusCode {
	case 400:
		return fmt.Errorf(`400: Bad Request - We couldn't fetch the resource you requested`)
	case 401:
		return fmt.Errorf(`401: Unauthorized - We couldn't fetch the resource you requested`)
	case 403:
		return fmt.Errorf(`403: Forbidden - We couldn't fetch the resource you requested`)
	case 404:
		return fmt.Errorf(`404: Not Found - We couldn't fetch the resource you requested`)
	case 500:
		return fmt.Errorf(`500: Internal Server Error - We couldn't fetch the resource you requested`)
	case 503:
		return fmt.Errorf(`503: Service Unavailable - We couldn't fetch the resource you requested`)
	case 504:
		return fmt.Errorf(`504: Gateway Timeout - We couldn't fetch the resource you requested`)
	}
	return httpErr
}
