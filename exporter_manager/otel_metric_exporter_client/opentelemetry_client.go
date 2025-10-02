package otelmetricexporterclient

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type ExporterClient struct {
	meterProvider *sdkmetric.MeterProvider
	destinations  []sdkmetric.Option
}

func (ec *ExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "otel_metric_exporter" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := ec.addOtelMetricExporter(ctx, exporter); err != nil {
				return err
			}
		}

		if exporter.GetDestinationType() == "splunk" {
			// Create OTLP gRPC exporter for Splunk collector
			if err := ec.addSplunkMetricExporter(ctx, exporter); err != nil {
				return err
			}
		}

		if exporter.GetDestinationType() == "honeycomb" {
			// Create OTLP gRPC exporter for Honeycomb collector
			if err := ec.addHoneycombMetricExporter(ctx, exporter); err != nil {
				return err
			}
		}
	}
	ec.meterProvider = sdkmetric.NewMeterProvider(ec.destinations...)

	return nil
}

func (ec *ExporterClient) addOtelMetricExporter(ctx context.Context, exporter e.Exporter) error {
	otlpExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	ec.destinations = append(ec.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(otlpExporter,
				sdkmetric.WithInterval(10*time.Second),
			)),
	)
	return nil
}

func (ec *ExporterClient) addSplunkMetricExporter(ctx context.Context, exporter e.Exporter) error {
	// 2. Splunk via OTLP HTTP (recommended approach)
	splunkExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("ingest.us0.signalfx.com"), // Change us0 to your realm
		otlpmetrichttp.WithURLPath("/v2/datapoint/otlp"),
		otlpmetrichttp.WithHeaders(map[string]string{
			"X-SF-Token": exporter.GetDestinationConfig().(api.SplunkConfig).Token,
		}),
		otlpmetrichttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		return err
	}

	ec.destinations = append(ec.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(splunkExporter,
				sdkmetric.WithInterval(10*time.Second),
			)),
	)
	return nil
}

func (ec *ExporterClient) addHoneycombMetricExporter(ctx context.Context, exporter e.Exporter) error {
	// 4. Additional OTLP HTTP exporter (e.g., for Honeycomb, New Relic, etc.)
	honeycombExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("api.honeycomb.io"),
		otlpmetrichttp.WithHeaders(map[string]string{
			"x-honeycomb-team": "YOUR_HONEYCOMB_API_KEY",
		}),
	)
	if err != nil {
		return err
	}

	ec.destinations = append(ec.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(honeycombExporter,
				sdkmetric.WithInterval(10*time.Second),
			)),
	)
	return nil
}
