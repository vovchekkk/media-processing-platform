package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	TaskDuration *prometheus.HistogramVec
	TasksTotal   *prometheus.CounterVec
}

func New() *Metrics {
	metrics := &Metrics{
		TaskDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "processor_task_duration_seconds",
				Help:    "Time spent processing tasks",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"filter"},
		),

		TasksTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "processor_tasks_total",
				Help: "Total number of tasks processed",
			},
			[]string{"filter", "status"},
		),
	}

	prometheus.MustRegister(metrics.TaskDuration, metrics.TasksTotal)

	return metrics
}