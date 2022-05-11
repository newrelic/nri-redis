package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
)

const transactionInContextKey = iota

const (
	appName                = "New Relic Redis Integration"
	apmInstrumentationName = "newrelic"
	waitForConnTimeout     = 5 * time.Second
)

type agentInstrumentationApm struct {
	nrApp     *newrelic.Application
	harvester *telemetry.Harvester
}

func NewAgentInstrumentationApm(license, apmEndpoint, telemetryEndpoint string) (AgentInstrumentation, error) {
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(c *newrelic.Config) {
			if apmEndpoint != "" {
				c.Host = apmEndpoint
			}
		},
		newrelic.ConfigDebugLogger(os.Stderr),
	)
	if err != nil {
		return nil, err
	}

	// Wait for the application to connect.
	if err = nrApp.WaitForConnection(waitForConnTimeout); nil != err {
		fmt.Println(err)
	}

	harvester, err := telemetry.NewHarvester(
		telemetry.ConfigAPIKey(license),
		func(c *telemetry.Config) {
			if telemetryEndpoint != "" {
				c.MetricsURLOverride = telemetryEndpoint
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &agentInstrumentationApm{nrApp: nrApp, harvester: harvester}, nil
}

func (a *agentInstrumentationApm) RecordMetric(ctx context.Context, metric metric) {
	var m telemetry.Metric
	switch metric.Type {
	case Gauge:
		m = telemetry.Gauge{
			Timestamp: metric.Timestamp, Value: metric.Value, Name: metric.Name, Attributes: metric.Attributes,
		}
	case Sum:
		m = telemetry.Count{
			Timestamp: metric.Timestamp, Value: metric.Value, Name: metric.Name, Attributes: metric.Attributes,
		}
	case Histrogram:
		//not implemented?
		return
	}
	a.harvester.RecordMetric(m)
}

func (a *agentInstrumentationApm) StartTransaction(ctx context.Context, name string) (context.Context, Transaction) {
	nrTxn := a.nrApp.StartTransaction(name)
	txn := &TransactionApm{nrTxn: nrTxn}
	ctx = ContextWithTransaction(ctx, txn)

	return ctx, txn
}

func (a *agentInstrumentationApm) Shutdown(timeout time.Duration) {
	a.nrApp.Shutdown(timeout)
}

type TransactionApm struct {
	nrTxn *newrelic.Transaction
}

func (t *TransactionApm) AddAttribute(key string, value interface{}) {
	t.nrTxn.AddAttribute(key, value)
}

func (t *TransactionApm) StartSegment(ctx context.Context, name string) (context.Context, Segment) {
	return ctx, t.nrTxn.StartSegment(name)
}

func (t *TransactionApm) StartExternalSegment(ctx context.Context, name string, req *http.Request) (context.Context, Segment) {
	return ctx, newrelic.StartExternalSegment(t.nrTxn, req)
}

func (t *TransactionApm) NoticeError(err error) {
	t.nrTxn.NoticeError(err)
}

func (t *TransactionApm) AcceptDistributedTraceHeaders(_ newrelic.TransportType, payload http.Header) {
	t.nrTxn.AcceptDistributedTraceHeaders(newrelic.TransportOther, payload)
}

func (t *TransactionApm) End() {
	t.nrTxn.End()
}

type SegmentApm struct {
	nrSeg *newrelic.Segment
	ctx   context.Context
}

func (t *SegmentApm) AddAttribute(key string, value interface{}) {
	t.nrSeg.AddAttribute(key, value)
}

func (t *SegmentApm) End() {
	t.nrSeg.End()
}
