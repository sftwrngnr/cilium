package metrics

// metrics holds prometheus metrics objects and related utility functions. It
// does not abstract away the prometheus client but the caller rarely needs to
// refer to prometheus directly.
//
// **Adding a metric**
// - Add a metric object of the appropriate type that is exported
// - Register the new object in the init function

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	registry = prometheus.NewRegistry()

	// Namespace is used to scope metrics from cilium. It is prepended to metric
	// names and separated with a '_'
	Namespace = "cilium"

	// Endpoints

	// NumEndpoints is the number of managed endpoints
	NumEndpoints = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "endpoints",
		Help:      "Number of endpoints managed by this agent",
		// FIXME: do we have any node IDs for the agent? or agent IDs? Do we even have to provide anything?
		// ConstLabels: prometheus.Labels{"node": "a node ID, from k8s?",
	})

	// NumEndpointsRegenerating is the number of endpoints currently regenerating
	NumEndpointsRegenerating = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "endpoints_regenerating",
		Help:      "Number of endpoints currently regenerating",
	})

	// CountEndpointsRegenerations is a count of the number of times any endpoint
	// has been regenerated and success/fail outcome
	CountEndpointsRegenerations = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "endpoints_regenerations",
		Help:      "Count of all endpoint regenerations that have completed, tagged by outcome",
	},
		[]string{"outcome"})

	// Policies

	// NumPolicies is the number of policies loaded into the agent
	NumPolicies = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "policies",
		Help:      "Number of policies currently loaded",
	})

	// PolicyRevision is the current policy revision number for this agent
	PolicyRevision = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "policies_max_revision",
		Help:      "Highest policy revision number in the agent",
	})

	// PolicyImportErrors is a count of failed policy imports
	PolicyImportErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "policies_import_errors",
		Help:      "Number of times a policy import has failed",
	})

	// Orchestration

	// LastK8sEventTS is the time in seconds since epoch that we last recieved an
	// event from k8s
	LastK8sEventTS = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "last_k8s_event_ts",
		Help:      "Last timestamp when we received a kubernetes event",
	})
)

func init() {
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), Namespace))
	// TODO: Figure out how to put this into a Namespace
	//registry.MustRegister(prometheus.NewGoCollector())

	registry.MustRegister(NumEndpoints)
	registry.MustRegister(NumEndpointsRegenerating)
	registry.MustRegister(CountEndpointsRegenerations)

	registry.MustRegister(NumPolicies)
	registry.MustRegister(PolicyRevision)
	registry.MustRegister(PolicyImportErrors)

	registry.MustRegister(LastK8sEventTS)
}

// Enable begins serving prometheus metrics on the address passed in. Addresses
// of the form ":8080" will bind the port on all interfaces.
func Enable(addr string) error {
	go func() {
		// The Handler function provides a default handler to expose metrics
		// via an HTTP server. "/metrics" is the usual endpoint for that.
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.WithError(http.ListenAndServe(addr, nil)).Warn("Cannot start metrics server on %s", addr)
	}()

	return nil
}

// SetTSValue sets the gauge to the time value provided
func SetTSValue(c prometheus.Gauge, ts time.Time) {
	// Build time in seconds since the epoch. Prometheus only takes floating
	// point values, however, and urges times to be in seconds
	c.Set(float64(ts.UnixNano()) / float64(1000000000))
}
