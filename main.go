package main

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	tracer := otel.Tracer("example/main")
	ctx, span := tracer.Start(ctx, "example")
	defer span.End()

	log.SetFormatter(&log.JSONFormatter{})

	standardFields := log.Fields{
		"dd.trace_id": convertTraceID(span.SpanContext().TraceID().String()),
		"dd.span_id":  convertTraceID(span.SpanContext().SpanID().String()),
		"dd.service":  "serviceName",
		"dd.env":      "serviceEnv",
		"dd.version":  "serviceVersion",
	}

	log.WithFields(standardFields).WithContext(ctx).Info("hello world")
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
