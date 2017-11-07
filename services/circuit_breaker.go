package services

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sony/gobreaker"
)

// CircuitBreaker - instance new circuit breaker
var CircuitBreaker = &MISCircuitBreaker{}

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	cb := gobreaker.NewCircuitBreaker(st)
	CircuitBreaker.CircuitBreaker = cb
}

// MISCircuitBreaker - cicrcuit breaker for MIS to external system
type MISCircuitBreaker struct {
	*gobreaker.CircuitBreaker
}

// Get - get http to external service
func (m *MISCircuitBreaker) Get(url string) ([]byte, error) {
	body, err := CircuitBreaker.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if resp == nil {
			return nil, errors.New("Service is unavailable")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	})
	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}
