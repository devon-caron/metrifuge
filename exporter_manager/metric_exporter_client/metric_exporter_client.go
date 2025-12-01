package metric_exporter_client

import (
	"context"
	"fmt"
	"time"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MetricExporterClient struct {
	meterProviders map[string]map[string]*sdkmetric.MeterProvider
	meters         map[string]map[string]metric.Meter
	destinations   []sdkmetric.Option
}

func (me *MetricExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "OtelCollector" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := me.addOtelCollector(ctx, exporter); err != nil {
				return fmt.Errorf("failed to add OTLP collector: %w", err)
			}
		} else if exporter.GetDestinationType() == "honeycomb" {
			// Create OTLP HTTP exporter for Honeycomb
			if err := me.addHoneycombMetricExporter(ctx, exporter); err != nil {
				return fmt.Errorf("failed to add Honeycomb exporter: %w", err)
			}
			// } else if exporter.GetDestinationType() == "prometheus" {
			// 	// Create OTLP gRPC exporter for Prometheus collector
			// 	if err := me.addPrometheusMetricExporter(ctx, exporter); err != nil {
			// 		return err
			// 	}
		} else {
			return fmt.Errorf("unknown destination type: %s", exporter.GetDestinationType())
		}
		if me.meterProviders == nil {
			me.meterProviders = make(map[string]map[string]*sdkmetric.MeterProvider)
		}
		if me.meters == nil {
			me.meters = make(map[string]map[string]metric.Meter)
		}
		ns := exporter.GetLogSourceInfo().Namespace
		name := exporter.GetLogSourceInfo().Name
		if me.meterProviders[ns] == nil {
			me.meterProviders[ns] = make(map[string]*sdkmetric.MeterProvider)
		}
		if me.meters[ns] == nil {
			me.meters[ns] = make(map[string]metric.Meter)
		}
		me.meterProviders[ns][name] = sdkmetric.NewMeterProvider(me.destinations...)
		me.meters[ns][name] = me.meterProviders[ns][name].Meter("metrifuge")
	}

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
	refreshInterval, err := time.ParseDuration(exporter.Spec.RefreshInterval)
	if err != nil {
		return fmt.Errorf("failed to parse refresh interval: %w", err)
	}
	me.destinations = append(me.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(otlpExporter,
				sdkmetric.WithInterval(refreshInterval),
			)),
	)
	return nil
}

func (me *MetricExporterClient) addHoneycombMetricExporter(ctx context.Context, exporter e.Exporter) error {
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
	honeycombExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("api.honeycomb.io"),
		otlpmetrichttp.WithHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("failed to create Honeycomb OTLP HTTP exporter: %w", err)
	}

	// Parse refresh interval
	refreshInterval, err := time.ParseDuration(exporter.Spec.RefreshInterval)
	if err != nil {
		return fmt.Errorf("failed to parse refresh interval: %w", err)
	}

	me.destinations = append(me.destinations,
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(honeycombExporter,
				sdkmetric.WithInterval(refreshInterval),
			)),
	)
	return nil
}

func (me *MetricExporterClient) ExportMetric(ctx context.Context, metricData *api.MetricData) error {
	// Get cached meter using the correct namespace and exporter name
	namespace := ""
	name := ""
	// Extract namespace and exporter name from context if available
	if ns, ok := ctx.Value(global.SOURCE_NAMESPACE_KEY).(string); ok && ns != "" {
		namespace = ns
	}
	if expName, ok := ctx.Value(global.SOURCE_NAME_KEY).(string); ok && expName != "" {
		name = expName
	}
	meter := me.meters[namespace][name]
	if meter == nil {
		return fmt.Errorf("meter not found for namespace %s and name %s", namespace, name)
	}

	// Create and record based on metric kind
	switch metricData.Kind {
	case "Int64Counter":
		counter, err := meter.Int64Counter(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create int64 counter: %w", err)
		}
		counter.Add(ctx, metricData.ValueInt, metric.WithAttributes(metricData.Attributes...))

	case "Float64Counter":
		counter, err := meter.Float64Counter(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create float64 counter: %w", err)
		}
		counter.Add(ctx, metricData.ValueFloat, metric.WithAttributes(metricData.Attributes...))

	case "Int64Gauge":
		gauge, err := meter.Int64Gauge(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create int64 gauge: %w", err)
		}
		gauge.Record(ctx, metricData.ValueInt, metric.WithAttributes(metricData.Attributes...))

	case "Float64Gauge":
		gauge, err := meter.Float64Gauge(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create float64 gauge: %w", err)
		}
		gauge.Record(ctx, metricData.ValueFloat, metric.WithAttributes(metricData.Attributes...))

	case "Int64Histogram":
		histogram, err := meter.Int64Histogram(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create int64 histogram: %w", err)
		}
		histogram.Record(ctx, metricData.ValueInt, metric.WithAttributes(metricData.Attributes...))

	case "Float64Histogram":
		histogram, err := meter.Float64Histogram(metricData.Name)
		if err != nil {
			return fmt.Errorf("failed to create float64 histogram: %w", err)
		}
		histogram.Record(ctx, metricData.ValueFloat, metric.WithAttributes(metricData.Attributes...))

	default:
		return fmt.Errorf("unsupported metric kind: %s", metricData.Kind)
	}

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
