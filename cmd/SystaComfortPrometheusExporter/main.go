//
// Main entry for the SystaComfort Prometheus exporter
//

package main

import (
    "flag"

    internal "github.com/xgcssch/SystaComfortPrometheusExporter/internal/pkg/udpserver"
    "k8s.io/klog/v2"
)

var registerGoCollector = flag.Bool("registerGoCollector", false, "Register the GO collector")
var registerProcessCollector = flag.Bool("registerProcessCollector", false, "Register the process collector")
var prometheusPort = flag.Int("port", 2112, "Port to use exposing the exporter")
var prometheusURL = flag.String("url", "/metrics", "URL where the metrics are exposed")

func main() {
    klog.InitFlags(nil)
    flag.Parse()

    klog.Info("SystaComfortPrometheusExporter v0.1 starting")

    Configuration := internal.ProgramConfiguration{
        PrometheusPort:  *prometheusPort,
        PrometheusURL: *prometheusURL,
        RegisterGoCollector: *registerGoCollector,
        RegisterProcessCollector: *registerProcessCollector,
    }

    internal.StartupServer( Configuration )
}
