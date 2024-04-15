package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func main() {
	projectID := flag.String("project", "", "Google Cloud Platform project ID")
	flag.Parse()
	if *projectID == "" {
		log.Println("project ID not set")
		os.Exit(-1)
	}

	resOpt := basic.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		attribute.String(metric.CloudKeyProvider, metric.CloudProviderGCP),
	))
	_, err := metric.InstallNewPipeline([]metric.Option{
		metric.WithProjectID(*projectID),
	}, resOpt)
	if err != nil {
		log.Printf("failed to install new pipeline: %v", err)
		os.Exit(-1)
	}

	meter, err := global.MeterProvider().Meter("my-latency-metric").
		SyncFloat64().Histogram("my-latency-metric",
		instrument.WithDescription("The distribution of the duration"),
		instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		log.Printf("could not create meter due to: %s", err)
		os.Exit(-1)
	}

	ctx := context.Background()
	startTime := time.Now()
	time.Sleep(time.Second)
	meter.Record(ctx, float64(time.Since(startTime).Milliseconds()), attribute.String("key", "value"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until a signal is received.
	<-c
}
