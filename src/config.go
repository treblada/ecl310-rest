package main

import "flag"

type CmdLineArgs struct {
	eclHost string
	eclPort int
}

func parseCmdLine() CmdLineArgs {
	host := flag.String("host", "localhost", "ECL310 hostname or IP address. Defaults to localhost")
	port := flag.Int("port", 502, "ECL310 MODbus port. Defaults to 502")
	flag.Parse()
	return CmdLineArgs{
		eclHost: *host,
		eclPort: *port,
	}
}
