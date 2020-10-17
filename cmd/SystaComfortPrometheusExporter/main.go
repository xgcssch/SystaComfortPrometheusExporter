//
// A perl script developed by Klaus.Schmidinger@tvdr.de

package main

import (
	internal "github.com/xgcssch/SystaComfortPrometheusExporter/internal/pkg/udpserver"
)

func main() {
	internal.StartupServer(22460)
}
