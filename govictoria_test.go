package govictoria

import (
	"testing"
	"time"
)

func Test_SendMetrics(t *testing.T) {
	vm := NewGoVictoria("http://localhost:8428", "", "")

	err := vm.SendMetrics([]VictoriaMetricsRequest{
		{
			Metric: map[string]string{
				"__name__": "test_metric",
			},
			Values:     []int64{1},
			Timestamps: []int64{time.Now().Unix()},
		},
	})

	if err != nil {
		t.Error(err)
	}
}

func Test_QueryTimeRange(t *testing.T) {
	vm := NewGoVictoria("http://localhost:8428", "", "")

	_, err := vm.QueryTimeRange("test_metric{}", time.Now().Add(-time.Hour*24), time.Now(), "20s")
	if err != nil {
		t.Error(err)
	}
}
