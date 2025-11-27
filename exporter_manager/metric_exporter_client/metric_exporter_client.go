package metric_exporter_client

import (
	"context"
	"fmt"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MetricExporterClient struct {
	meterProvider *sdkmetric.MeterProvider
	destinations  []sdkmetric.Option
}

func (me *MetricExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "OtelCollector" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := me.addOtelCollector(ctx, exporter); err != nil {
				return err
			}
			// } else if exporter.GetDestinationType() == "honeycomb" {
			// 	// Create OTLP gRPC exporter for Honeycomb collector
			// 	if err := me.addHoneycombMetricExporter(ctx, exporter); err != nil {
			// 		return err
			// 	}
			// } else if exporter.GetDestinationType() == "prometheus" {
			// 	// Create OTLP gRPC exporter for Prometheus collector
			// 	if err := me.addPrometheusMetricExporter(ctx, exporter); err != nil {
			// 		return err
			// 	}
		} else {
			return fmt.Errorf("unknown destination type: %s", exporter.GetDestinationType())
		}
	}
	me.meterProvider = sdkmetric.NewMeterProvider(me.destinations...)

	return nil
}

func (me *MetricExporterClient) addOtelCollector(ctx context.Context, exporter e.Exporter) error {

	endpoint := exporter.Spec.Destination.OtelCollector.Endpoint
	if endpoint == "" {
		return fmt.Errorf("otel collector endpoint is required")
	}
	insecure := exporter.Spec.Destination.OtelCollector.Insecure

	options := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
	if insecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	otlpExporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
	}
	me.destinations = append(me.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(otlpExporter,
				sdkmetric.WithInterval(10*time.Second),
			)),
	)
	return nil
}

func (me *MetricExporterClient) addHoneycombMetricExporter(ctx context.Context, exporter e.Exporter) error {
	// 4. Additional OTLP HTTP exporter (e.g., for Honeycomb, New Relic, etc.)
	honeycombExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("api.honeycomb.io"),
		otlpmetrichttp.WithHeaders(map[string]string{
			"x-honeycomb-team": "YOUR_HONEYCOMB_API_KEY",
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create Honeycomb OTLP HTTP exporter: %w", err)
	}

	me.destinations = append(me.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(honeycombExporter,
				sdkmetric.WithInterval(10*time.Second),
			)),
	)
	return nil
}

func (me *MetricExporterClient) ExportMetric(ctx context.Context, metric *api.MetricData) error {
	// TODO: Implement actual metric export logic
	return nil
}

// func (me *MetricExporterClient) addPrometheusMetricExporter(ctx context.Context, exporter e.Exporter) error {
// 	prometheusExporter, err := prometheus.New()
// 	if err != nil {
// 		return err
// 	}

// 	me.destinations = append(me.destinations,
// 		sdkmetric.WithReader(
// 			sdkmetric.NewPeriodicReader(prometheusExporter,
// 				sdkmetric.WithInterval(10*time.Second),
// 			)),
// 	)
// 	return nil
// }
