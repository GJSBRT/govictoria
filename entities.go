package govictoria

type GoVictoriaConfig struct {
	URL      string
	Username string
	Password string
}

type VictoriaMetricsRequest struct {
	Metric     map[string]string `json:"metric"`
	Values     []int64           `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

type VictoriaMetricsQueryResponse struct {
	Status string `json:"status"`
	Data  struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values []int   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
