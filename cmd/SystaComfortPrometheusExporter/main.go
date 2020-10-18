//
// A perl script developed by Klaus.Schmidinger@tvdr.de
// go run github.com/xgcssch/SystaComfortPrometheusExporter/cmd/SystaComfortPrometheusExporter

package main

import (
	internal "github.com/xgcssch/SystaComfortPrometheusExporter/internal/pkg/udpserver"
)

func main() {
	internal.StartupServer(22460)
}
