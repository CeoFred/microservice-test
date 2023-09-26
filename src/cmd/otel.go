package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, serviceName, serviceVersion string, cfg *config.Config) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Setup resource.
	res, err := newResource(serviceName, serviceVersion, cfg.Environment)
	if err != nil {
		handleErr(err)
		return
	}

	if cfg.TraceDestination != "" {
		// Setup trace provider.

		f, err := os.Create(cfg.TraceDestination)
		if err != nil {
			handleErr(err)
			return nil, err
		}
		shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
			return f.Close()
		})

		tracerProvider, err := newTraceProvider(res, f)
		if err != nil {
			handleErr(err)
			return nil, err
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	} else {
		log.Info().Msg("Skipping trace logging")
	}

	if cfg.MetricsDestination != "" {
		// Setup meter provider.

		f, err := os.Create(cfg.MetricsDestination)
		if err != nil {
			handleErr(err)
			return nil, err
		}
		shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
			return f.Close()
		})

		meterProvider, err := newMeterProvider(res, f)
		if err != nil {
			handleErr(err)
			return nil, err
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)

	} else {
		log.Info().Msg("Skipping metrics logging")
	}

	return
}

func newResource(serviceName, serviceVersion, serviceEnvironment string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			attribute.String("environment", serviceEnvironment),
		))
}

func newTraceProvider(res *resource.Resource, w io.Writer) (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(w),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(5*time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newMeterProvider(res *resource.Resource, w io.Writer) (*metric.MeterProvider, error) {
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetIndent("", "  ")

	metricExporter, err := stdoutmetric.New(
		stdoutmetric.WithEncoder(jsonEncoder),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(1*time.Minute))),
	)
	return meterProvider, nil
}
