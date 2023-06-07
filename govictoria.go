package govictoria

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

	err = response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Victoria Metrics returned a non-200 status code")
	}

	return nil
}
