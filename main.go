package main

import (
	"flag"

	"github.com/SebastianCzoch/lxc-exporter/service"
)

var addr = flag.String("web.listen-address", ":9125", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	service.StartColectors()
	service.StartServer(addr)
}
