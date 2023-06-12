package main

import (
	"time"

	"github.com/betelgeusexru/golang-toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

type MetricsMiddleware struct { 
	errCounterAgg prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name: "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name: "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name: "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name: "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name: "aggregate",
		Buckets: []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name: "calculate",
		Buckets: []float64{0.1, 0.5, 1},
	})
	
	return &MetricsMiddleware{
		next: next,
		errCounterAgg: errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg: reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg: reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time){
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return err
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time){
		m.reqLatencyCalc.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	
	return m.next.CalculateInvoice(obuID)
}


func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time){
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
		}).Info("AggregateDistance")
	}(time.Now())
	return m.next.AggregateDistance(distance)
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time){
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
		}).Info("CalculateInvoice")
	}(time.Now())
	return m.next.CalculateInvoice(obuID)
}