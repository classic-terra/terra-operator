package main

import (
	"flag"
)

var (
	networkReplica int
	configDir      string
)

func GetNetworkFlag() *flag.FlagSet {
	networkFlags := flag.NewFlagSet("network", flag.ExitOnError)
	networkFlags.IntVar(&networkReplica, "replica", 4, "Number of network replica")
	networkFlags.StringVar(&configDir, "config-dir", "network_deployment.yaml", "Config file directory for network")

	return networkFlags
}
