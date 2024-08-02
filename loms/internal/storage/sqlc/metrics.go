package sqlc

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
	"time"
)

var sqlQueryCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app",
		Name:      "sql_query_total_counter",
		Help:      "Total amount of sql queries",
	},
	[]string{"type", "status"},
)

var sqlQueryHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app",
		Name:      "sql_query_duration_histogram",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"type", "status"})

var queryTypes = map[string]string{
	"select": "select",
	"insert": "insert",
	"update": "update",
	"delete": "delete",
}

func recordMetrics(sql string, startTime time.Time, err error) {
	if len(sql) < 6 {
		return
	}

	sc := bufio.NewScanner(strings.NewReader(sql))
	sc.Scan()
	sc.Scan()
	rawSql := sc.Text()

	testWord := strings.ToLower(rawSql[:6])

	queryType, ok := queryTypes[testWord]
	if !ok {
		return
	}

	status := "success"
	if err != nil {
		status = "error"
	}

	sqlQueryCounterVec.WithLabelValues(queryType, status).Inc()
	sqlQueryHistogramVec.WithLabelValues(queryType, status).Observe(time.Since(startTime).Seconds())

}
