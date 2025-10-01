package otelmetricexporterclient

import (
	"context"

	"github.com/devon-caron/metrifuge/k8s/api"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type ExporterClient struct {
	meterProvider *sdkmetric.MeterProvider
}

func (ec *ExporterClient) Initialize(ctx context.Context, exporters []api.Exporter) error {
	var destinations = make([]sdkmetric.Option, 0)
	var otelAddedFlag bool = false
	for _, exporter := range exporters {
		if exporter.GetDestinationType() == "otel_metric_exporter" && !otelAddedFlag {
			// Create OTLP gRPC exporter for OTEL collector
			otlpExporter, err := otlpmetricgrpc.New(ctx,
				otlpmetricgrpc.WithEndpoint("localhost:4317"),
				otlpmetricgrpc.WithInsecure(),
			)
			if err != nil {
				return err
			}
			destinations = append(destinations, sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(otlpExporter),
			))
		}
	}
	ec.meterProvider = sdkmetric.NewMeterProvider(destinations...)

	return nil
}
