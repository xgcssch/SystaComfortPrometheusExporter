//
// Main entry for the SystaComfort Prometheus exporter
//

package main

import (
	"flag"

	internal "github.com/xgcssch/SystaComfortPrometheusExporter/internal/pkg/udpserver"
)

var prometheusPort = flag.Int("port", 2112, "Port to use exposing the exporter")
var prometheusURL = flag.String("url", "/metrics", "URL where the metrics are exposed")
var dumpValues = flag.Bool("dump", false, "Dump values received from heating controller")

func main() {
	flag.Parse()

	internal.StartupServer(
		*prometheusPort,
		*prometheusURL,
		*dumpValues)
}
