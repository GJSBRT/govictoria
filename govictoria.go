package govictoria

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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
	body := ""

	// Loop through the request and build the body
	for _, requestBody := range requests {
		jsonRequest, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		body += string(jsonRequest)
	}

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("POST", g.Config.URL+"/api/v1/import", bytes.NewBuffer([]byte(body)))
	request.Header.Add("Authorization", "Basic "+basicAuth(g.Config.Username, g.Config.Password))

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
		return errors.New("Victoria Metrics returned a non-200 status code")
	}

	return nil
}

// QueryTimeRange queries Victoria Metrics for metrics in a time range
func (g *GoVictoria) QueryTimeRange(promql string, startTime time.Time, endTime time.Time, step string) (metrics map[string]string, err error) {
	// Check if the start time is before the end time
	if startTime.After(endTime) {
		return nil, errors.New("Start time must be before end time")
	}

	// Add the query parameters to the request
	url := g.Config.URL
	url += "/api/v1/query_range"
	url += "?query=" + promql
	url += "&start=" + strconv.FormatInt(startTime.Unix(), 10)
	url += "&end=" + strconv.FormatInt(endTime.Unix(), 10)
	url += "&step=" + step

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("GET", url, nil)

	// Add the query parameters to the request
	request.Header.Add("Authorization", "Basic "+basicAuth(g.Config.Username, g.Config.Password))

	// Send the request to Victoria Metrics
	response, err := g.Client.Do(request)
	if err != nil {
		return nil, err
	}

	// Close the response body
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	// Check if the status code is not 200
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Victoria Metrics returned a non-200 status code")
	}

	return nil, nil
}
