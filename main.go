package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "Polygon"
	environment = "production"
	id          = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(

		tracesdk.WithBatcher(exp),

		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {

		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tracer := otel.Tracer("component-main")

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	allFunctions(ctx)
}

func allFunctions(ctx context.Context) {

	var n uint = 15

	ansFib := FibFunc(ctx, n)
	ansFact := FactFunc(ctx, n)

	fmt.Print("\nCalculating for 15 :\n")
	fmt.Print("\n Fibbonaci :", ansFib)
	fmt.Print("\n Factorial :", ansFact, "\n")
}

func FibFunc(ctx context.Context, n uint) uint {

	tracer := otel.Tracer("component fibonacci")
	_, span := tracer.Start(ctx, "fibbonaci")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	ansFib, fibErr := Fibonacci(n)

	time.Sleep(2 * time.Second) //Adding 2 secs sleep

	if fibErr != nil {
		panic((fibErr))
	}
	return uint(ansFib)
}

func FactFunc(ctx context.Context, n uint) uint {

	tracer := otel.Tracer("component factorial")
	_, span := tracer.Start(ctx, "factorial")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	ansFact := Factorial(int(n))

	time.Sleep(1 * time.Second) //Adding 1 sec sleep

	return (uint(ansFact))
}
