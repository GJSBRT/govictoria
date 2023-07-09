package govictoria

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GoVictoria is the main struct for the library
type GoVictoria struct {
	// Config is the configuration for the library
	Config *GoVictoriaConfig

	// The HTTP client to use for requests
	Client *http.Client
}

// NewGoVictoria creates a new GoVictoria instance
func NewGoVictoria(url string, username string, password string) *GoVictoria {
	return &GoVictoria{
		Config: &GoVictoriaConfig{
			URL:      url,
			Username: username,
			Password: password,
		},
		Client: &http.Client{},
	}
}

// SendMetrics sends the metrics to VictoriaMetrics
func (g *GoVictoria) SendMetrics(requests []VictoriaMetricsRequest) error {
	if len(requests) == 0 {
		return errors.New("No requests to send")
	}

	// Loop through the request and build the body
	body := ""
	for _, requestBody := range requests {
		jsonRequest, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		body += string(jsonRequest)
	}

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("POST", g.Config.URL+"/api/v1/import", bytes.NewBuffer([]byte(body)))
	request.Header.Add("Authorization", "Basic "+BasicAuth(g.Config.Username, g.Config.Password))

	// Send the request to Victoria Metrics
	response, err := g.Client.Do(request)
	if err != nil {
		return err
	}

	// Close the response body
	err = response.Body.Close()
	if err != nil {
		return err
	}

	// Check if the status code is not 204
	if response.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Victoria Metrics returned a non-200 status code: %d", response.StatusCode))
	}

	return nil
}

// QueryTimeRange queries Victoria Metrics for metrics in a time range
func (g *GoVictoria) QueryTimeRange(promql string, startTime time.Time, endTime time.Time, step string) (VictoriaMetricsQueryResponse, error) {
	// Check if the start time is before the end time
	if startTime.After(endTime) {
		return VictoriaMetricsQueryResponse{}, errors.New("Start time must be before end time")
	}

	// Add the query parameters to the request
	params := url.Values{}
	params.Add("query", promql)
	params.Add("start", strconv.FormatInt(startTime.Unix(), 10))
	params.Add("end", strconv.FormatInt(endTime.Unix(), 10))
	params.Add("step", step)

	url := g.Config.URL + "/api/v1/query_range?" + params.Encode()

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Add the query parameters to the request
	request.Header.Add("Authorization", "Basic "+BasicAuth(g.Config.Username, g.Config.Password))

	// Send the request to Victoria Metrics
	response, err := g.Client.Do(request)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Check if the status code is not 200
	if response.StatusCode != http.StatusOK {
		return VictoriaMetricsQueryResponse{}, errors.New(fmt.Sprintf("Victoria Metrics returned a non-200 status code: %d", response.StatusCode))
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Close the response body
	err = response.Body.Close()
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Unmarshal the response
	var metrics VictoriaMetricsQueryResponse
	err = json.Unmarshal([]byte(body), &metrics)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	return metrics, nil
}
