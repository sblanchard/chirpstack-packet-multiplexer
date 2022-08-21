package main

import "github.com/sblanchard/chirpstack-packet-multiplexer/cmd/chirpstack-packet-multiplexer/cmd"

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
