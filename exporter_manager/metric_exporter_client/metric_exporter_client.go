package metric_exporter_client

import (
	"context"
	"fmt"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MetricExporterClient struct {
	meterProviders map[string]map[string]*sdkmetric.MeterProvider
	destinations   []sdkmetric.Option
}

func (me *MetricExporterClient) Initialize(ctx context.Context, exporters []e.Exporter) error {
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "OtelCollector" {
			// Create OTLP gRPC exporter for OTEL collector
			if err := me.addOtelCollector(ctx, exporter); err != nil {
				return fmt.Errorf("failed to add OTLP collector: %w", err)
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
		if me.meterProviders == nil {
			me.meterProviders = make(map[string]map[string]*sdkmetric.MeterProvider)
		}
		ns := exporter.GetLogSourceInfo().Namespace
		if me.meterProviders[ns] == nil {
			me.meterProviders[ns] = make(map[string]*sdkmetric.MeterProvider)
		}
		me.meterProviders[ns][exporter.GetLogSourceInfo().Name] = sdkmetric.NewMeterProvider(me.destinations...)
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

func (me *MetricExporterClient) ExportMetric(ctx context.Context, metricData *api.MetricData) error {
	// TODO impl rpoperly; Get a meter from the provider using the correct namespace and exporter name
	// For now, using a default exporter name - this should be improved
	meter := me.meterProviders["default"]["default"].Meter("metrifuge")

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
