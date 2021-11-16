package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// rules := []tracer.SamplingRule{tracer.RateRule(1)}

	f, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	tracer.Start(
		// tracer.WithSampling(rules),
		tracer.WithService("test-service"),
		tracer.WithEnv("test-env"),
	)
	defer tracer.Stop()

	if err := profiler.Start(
		profiler.WithService("test-service"),
		profiler.WithEnv("test-env"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	); err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	// Create a traced mux router.
	mux := muxtrace.NewRouter()
	// Continue using the router as you normally would.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		span := tracer.StartSpan("web.request", tracer.ResourceName("parent"))
		defer span.Finish()

		child1 := tracer.StartSpan("child1", tracer.ResourceName("fibbonaci"), tracer.ChildOf(span.Context()))

		fibAnswer, fibErr := Fibonacci(12)

		fmt.Print("\n", fibAnswer, "\n", fibErr)
		time.Sleep(1 * time.Second)

		child1.Finish()

		child2 := tracer.StartSpan("child2", tracer.ResourceName("factorial"), tracer.ChildOf(span.Context()))

		factAnswer := Factorial(12)

		fmt.Print("\n", factAnswer)
		time.Sleep(2 * time.Second)

		child2.Finish()

		w.Write([]byte("Hello World!"))

		log.WithFields(log.Fields{
			"Fibonacci":   fibAnswer,
			"Factorial":   factAnswer,
			"dd.trace_id": child1.Context().TraceID(),
			"dd.span_id":  child1.Context().SpanID(),
		}).Info("Input Number is 12")

		// log.Printf("CHILD 1 Log Message : %v", child1)

	})
	http.ListenAndServe(":8080", mux)
}
