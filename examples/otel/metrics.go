package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc/credentials"
)

func setupMetrics(ctx context.Context, serviceName string) (*metric.MeterProvider, error) {
	c, err := getTls()
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithTLSCredentials(
			// mutual tls.
			credentials.NewTLS(c),
		),
	)
	if err != nil {
		return nil, err
	}

	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("metrics-attribute", "from-metrics"),
	)

	mp := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			// collects and exports metric data every 30 seconds.
			metric.NewPeriodicReader(exporter, metric.WithInterval(30*time.Second)),
		),
	)

	global.SetMeterProvider(mp)

	return mp, nil
}
