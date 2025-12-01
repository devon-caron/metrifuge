package log_exporter_client

import (
	"context"
	"fmt"
	"time"

	"github.com/devon-caron/metrifuge/global"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type LogExporterClient struct {
	loggerProviders map[string]map[string]*sdklog.LoggerProvider
	loggers         map[string]map[string]log.Logger
	destinations    map[string]map[string][]sdklog.LoggerProviderOption
}

func (le *LogExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "OtelCollector" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := le.addOtelCollector(ctx, exporter); err != nil {
				return fmt.Errorf("failed to add Otel collector: %v", err)
			}
		} else if exporter.GetDestinationType() == "honeycomb" {
			// Create OTLP HTTP exporter for Honeycomb
			if err := le.addHoneycombLogExporter(ctx, exporter); err != nil {
				return fmt.Errorf("failed to add Honeycomb log exporter: %w", err)
			}
		} else {
			return fmt.Errorf("unknown destination type: %s", exporter.GetDestinationType())
		}
	}

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
	if le.loggerProviders == nil {
		le.loggerProviders = make(map[string]map[string]*sdklog.LoggerProvider)
	}
	if le.loggers == nil {
		le.loggers = make(map[string]map[string]log.Logger)
	}
	if le.destinations == nil {
		le.destinations = make(map[string]map[string][]sdklog.LoggerProviderOption)
	}
	ns := exporter.GetLogSourceInfo().Namespace
	name := exporter.GetLogSourceInfo().Name
	if le.loggerProviders[ns] == nil {
		le.loggerProviders[ns] = make(map[string]*sdklog.LoggerProvider)
	}
	if le.loggers[ns] == nil {
		le.loggers[ns] = make(map[string]log.Logger)
	}
	if le.destinations[ns] == nil {
		le.destinations[ns] = make(map[string][]sdklog.LoggerProviderOption)
	}
	le.destinations[ns][name] = append(le.destinations[ns][name],
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(otlpExporter),
		),
	)
	le.loggerProviders[ns][name] = sdklog.NewLoggerProvider(le.destinations[ns][name]...)
	le.loggers[ns][name] = le.loggerProviders[ns][name].Logger("metrifuge")
	return nil
}

func (le *LogExporterClient) addHoneycombLogExporter(ctx context.Context, exporter e.Exporter) error {
	// Validate Honeycomb config
	honeycombConfig := exporter.Spec.Destination.Honeycomb
	if honeycombConfig == nil {
		return fmt.Errorf("honeycomb configuration is required")
	}
	if honeycombConfig.APIKey == "" {
		return fmt.Errorf("honeycomb API key is required")
	}
	if honeycombConfig.Dataset == "" {
		return fmt.Errorf("honeycomb dataset is required")
	}

	// Build headers for Honeycomb
	headers := map[string]string{
		"x-honeycomb-team":    honeycombConfig.APIKey,
		"x-honeycomb-dataset": honeycombConfig.Dataset,
	}

	// Add environment header if specified
	if honeycombConfig.Environment != "" {
		headers["x-honeycomb-environment"] = honeycombConfig.Environment
	}

	// Create OTLP HTTP exporter for Honeycomb
	honeycombExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("api.honeycomb.io"),
		otlploghttp.WithHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("failed to create Honeycomb OTLP HTTP log exporter: %w", err)
	}

	// Initialize maps if needed
	if le.loggerProviders == nil {
		le.loggerProviders = make(map[string]map[string]*sdklog.LoggerProvider)
	}
	if le.loggers == nil {
		le.loggers = make(map[string]map[string]log.Logger)
	}
	if le.destinations == nil {
		le.destinations = make(map[string]map[string][]sdklog.LoggerProviderOption)
	}

	ns := exporter.GetLogSourceInfo().Namespace
	name := exporter.GetLogSourceInfo().Name

	if le.loggerProviders[ns] == nil {
		le.loggerProviders[ns] = make(map[string]*sdklog.LoggerProvider)
	}
	if le.loggers[ns] == nil {
		le.loggers[ns] = make(map[string]log.Logger)
	}
	if le.destinations[ns] == nil {
		le.destinations[ns] = make(map[string][]sdklog.LoggerProviderOption)
	}

	le.destinations[ns][name] = append(le.destinations[ns][name],
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(honeycombExporter),
		),
	)
	le.loggerProviders[ns][name] = sdklog.NewLoggerProvider(le.destinations[ns][name]...)
	le.loggers[ns][name] = le.loggerProviders[ns][name].Logger("metrifuge")
	return nil
}

func (le *LogExporterClient) ExportLog(ctx context.Context, logMessage string) error {
	// Get cached logger using the correct namespace and exporter name
	namespace := ""
	name := ""
	// Extract namespace and exporter name from context if available
	if ns, ok := ctx.Value(global.SOURCE_NAMESPACE_KEY).(string); ok && ns != "" {
		namespace = ns
	}
	if expName, ok := ctx.Value(global.SOURCE_NAME_KEY).(string); ok && expName != "" {
		name = expName
	}
	logger := le.loggers[namespace][name]
	if logger == nil {
		return fmt.Errorf("logger not found for namespace %s and name %s", namespace, name)
	}

	// Create a log record
	var record log.Record
	record.SetTimestamp(time.Now())
	record.SetBody(log.StringValue(logMessage))
	record.SetSeverity(log.SeverityInfo)

	// Emit the log record
	logger.Emit(ctx, record)

	return nil
}
