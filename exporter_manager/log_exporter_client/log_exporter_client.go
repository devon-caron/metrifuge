package log_exporter_client

import (
	"context"
	"fmt"

	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type LogExporterClient struct {
	loggerProvider *sdklog.LoggerProvider
	destinations   []sdklog.LoggerProviderOption
}

func (le *LogExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "OtelCollector" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := le.addOtelCollector(ctx, exporter); err != nil {
				return err
			}
			// } else if exporter.GetDestinationType() == "honeycomb" {
			// 	// Create OTLP HTTP exporter for Honeycomb collector
			// 	if err := le.addHoneycombLogExporter(ctx, exporter); err != nil {
			// 		return err
			// 	}
		} else {
			return fmt.Errorf("unknown destination type: %s", exporter.GetDestinationType())
		}
	}
	le.loggerProvider = sdklog.NewLoggerProvider(le.destinations...)

	return nil
}

func (le *LogExporterClient) addOtelCollector(ctx context.Context, exporter e.Exporter) error {

	endpoint := exporter.Spec.Destination.OtelCollector.Endpoint
	if endpoint == "" {
		return fmt.Errorf("otel collector endpoint is required")
	}
	insecure := exporter.Spec.Destination.OtelCollector.Insecure

	options := []otlploggrpc.Option{otlploggrpc.WithEndpoint(endpoint)}
	if insecure {
		options = append(options, otlploggrpc.WithInsecure())
	}

	otlpExporter, err := otlploggrpc.New(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP gRPC log exporter: %w", err)
	}
	le.destinations = append(le.destinations,
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(otlpExporter),
		),
	)
	return nil
}

func (le *LogExporterClient) addHoneycombLogExporter(ctx context.Context, exporter e.Exporter) error {
	// OTLP HTTP exporter (e.g., for Honeycomb, New Relic, etc.)
	honeycombExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("api.honeycomb.io"),
		otlploghttp.WithHeaders(map[string]string{
			"x-honeycomb-team": "YOUR_HONEYCOMB_API_KEY",
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create Honeycomb OTLP HTTP log exporter: %w", err)
	}

	le.destinations = append(le.destinations,
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(honeycombExporter),
		),
	)
	return nil
}

func (le *LogExporterClient) ExportLog(ctx context.Context, logMessage string) error {
	// TODO: Implement actual log export logic
	return nil
}
