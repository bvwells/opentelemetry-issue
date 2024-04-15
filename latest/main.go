package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	otelsdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func main() {
	projectID := flag.String("project", "", "Google Cloud Platform project ID")
	flag.Parse()
	if *projectID == "" {
		log.Println("project ID not set")
		os.Exit(-1)
	}

	exporter, err := metric.New(
		metric.WithProjectID(*projectID),
		metric.WithMetricDescriptorTypeFormatter(func(m metricdata.Metrics) string {
			return fmt.Sprintf("custom.googleapis.com/opentelemetry/%s", m.Name)
		}))
	if err != nil {
		log.Printf("failed to create the exporter: %v", err)
		os.Exit(-1)
	}

	otel.SetMeterProvider(otelsdkmetric.NewMeterProvider(
		otelsdkmetric.WithReader(otelsdkmetric.NewPeriodicReader(exporter)),
	))

	provider := otel.GetMeterProvider()
	meter, err := provider.Meter(
		"my-latency-metric").
		Float64Histogram("my-latency-metric",
			otelmetric.WithDescription("The distribution of the duration"),
			otelmetric.WithUnit("ms"),
		)
	if err != nil {
		log.Printf("could not create meter due to: %v\n", err)
		os.Exit(-1)
	}

	ctx := context.Background()
	startTime := time.Now()
	time.Sleep(time.Second)
	meter.Record(ctx, float64(time.Since(startTime).Milliseconds()), otelmetric.WithAttributes(attribute.String("key", "value")))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until a signal is received.
	<-c
}
