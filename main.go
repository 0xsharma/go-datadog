package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// rules := []tracer.SamplingRule{tracer.RateRule(1)}

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

			// The profiles below are disabled by
			// default to keep overhead low, but
			// can be enabled as needed.
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
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
	})
	http.ListenAndServe(":8080", mux)
}
