package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func main() {
	for i := 0; i < 3; i++ {
		s := genStr(i)
		log.Println(s)
	}

	var (
		reg           = prometheus.NewRegistry()
		factory       = promauto.With(reg)
		randomNumbers = factory.NewHistogram(prometheus.HistogramOpts{
			Name:    "random_numbers",
			Help:    "A histogram of normally distributed random numbers.",
			Buckets: prometheus.LinearBuckets(-3, .1, 61),
		})
		requestCount = factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests by status code and method.",
			},
			[]string{"code", "method"},
		)
	)

	log.Println(randomNumbers != nil && requestCount != nil)
}

func genStr(i int) string {
	return fmt.Sprintf("%c, for %d", 'a'+i%26, i)
}
