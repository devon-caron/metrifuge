package otellogexporterclient

import (
	"context"
	"fmt"
	"time"

	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type LogExporterClient struct {
	loggerProvider *sdklog.LoggerProvider
	exporters      []sdklog.Exporter
}

func (le *LogExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "otel_log_exporter" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := le.addOtelLogExporter(ctx, exporter); err != nil {
				return err
			}
		} else if exporter.GetDestinationType() == "honeycomb_logs" {
			// Create OTLP HTTP exporter for Honeycomb logs
			if err := le.addHoneycombLogExporter(ctx, exporter); err != nil {
				return err
			}
		}
	}

	// Create logger provider with batch processor
	if len(le.exporters) > 0 {
		processors := make([]sdklog.LogRecordProcessor, 0, len(le.exporters))
		for _, exp := range le.exporters {
			processors = append(processors, sdklog.NewBatchProcessor(exp))
		}
		le.loggerProvider = sdklog.NewLoggerProvider(
			sdklog.WithProcessor(processors[0]), // Primary processor
		)
	}

	return nil
}

func (le *LogExporterClient) addOtelLogExporter(ctx context.Context, exporter e.Exporter) error {
	otlpExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint("localhost:4317"),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}
	le.exporters = append(le.exporters, otlpExporter)
	return nil
}

func (le *LogExporterClient) addHoneycombLogExporter(ctx context.Context, exporter e.Exporter) error {
	honeycombExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("api.honeycomb.io"),
		otlploghttp.WithHeaders(map[string]string{
			"x-honeycomb-team": "YOUR_HONEYCOMB_API_KEY",
		}),
		otlploghttp.WithTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create Honeycomb log exporter: %w", err)
	}
	le.exporters = append(le.exporters, honeycombExporter)
	return nil
}

// GetLoggerProvider returns the logger provider for use in the application
func (le *LogExporterClient) GetLoggerProvider() *sdklog.LoggerProvider {
	return le.loggerProvider
}

// EmitLog sends a log record to the configured exporters
func (le *LogExporterClient) EmitLog(ctx context.Context, timestamp time.Time, severity log.Severity, body string, attributes map[string]interface{}) error {
	if le.loggerProvider == nil {
		return fmt.Errorf("logger provider not initialized")
	}

	logger := le.loggerProvider.Logger("metrifuge")

	// Create log record
	record := sdklog.Record{}
	record.SetTimestamp(timestamp)
	record.SetSeverity(severity)
	record.SetBody(log.StringValue(body))

	// Add attributes
	for key, value := range attributes {
		switch v := value.(type) {
		case string:
			record.AddAttributes(log.String(key, v))
		case int:
			record.AddAttributes(log.Int(key, v))
		case int64:
			record.AddAttributes(log.Int64(key, v))
		case float64:
			record.AddAttributes(log.Float64(key, v))
		case bool:
			record.AddAttributes(log.Bool(key, v))
		default:
			record.AddAttributes(log.String(key, fmt.Sprintf("%v", v)))
		}
	}

	// Emit the log
	logger.Emit(ctx, record)
	return nil
}

// Shutdown gracefully shuts down the logger provider
func (le *LogExporterClient) Shutdown(ctx context.Context) error {
	if le.loggerProvider != nil {
		return le.loggerProvider.Shutdown(ctx)
	}
	return nil
}
