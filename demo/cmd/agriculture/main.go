package main

import (
	"flag"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
)

type metric struct {
	metric   string
	location string
	min      float64
	max      float64
	current  float64
	step     float64
}

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

var (
	gaugeMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "agriculture_metrics",
			Help: "agriculture metrics",
		},
		[]string{"metric", "location"},
	)

	alerts = []*metric{
		{"temperature", "雄安", 15, 25, 15, 0.1},
		{"temperature", "万全", 10, 20, 20, 0.07},
		{"temperature", "官厅湖", 20, 35, 20, 0.05},

		{"humidity", "雄安", 60, 80, 70, 0.2},
		{"humidity", "万全", 40, 50, 40, 0.1},
		{"humidity", "官厅湖", 30, 35, 30, 0.08},

		{"pm2.5", "雄安", 10, 30, 15, 0.2},
		{"pm2.5", "万全", 20, 40, 25, 0.1},
		{"pm2.5", "官厅湖", 30, 45, 35, 0.1},

		{"atmosPressure", "雄安", 103, 110, 103, 0.02},
		{"atmosPressure", "万全", 100, 102, 100, 0.01},
		{"atmosPressure", "官厅湖", 98, 100, 99, 0.01},
	}
)

func main() {

	flag.Parse()

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {

			for _, alert := range alerts {
				gaugeMetrics.WithLabelValues(alert.metric, alert.location).Set(alert.current)
				if alert.current > alert.max || alert.current < alert.min {
					alert.step = -alert.step
				}
				alert.current += alert.step
			}
		}
	}()

	log.Println("Listening on: ", *addr)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))

}
