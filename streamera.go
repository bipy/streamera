package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"streamera/client"
	"streamera/common"
	"streamera/server"
)

var mode string

func main() {
	if mode == "" || (mode != "server" && mode != "client") {
		usage()
		return
	}
	defer fmt.Println(common.Green("Bye"))
	if mode == "server" {
		serverModeInit()
	} else {
		clientModeInit()
	}
}

func init() {
	flag.Usage = usage
	flag.StringVar(&mode, "m", "", "run mode: \"server\" or \"client\"")
	flag.Parse()
}

func clientModeInit() {
	config, err := getConfig("client_config.json")
	if err != nil {
		panic(err)
	}

	c, err := client.NewClient(net.ParseIP(config.Ip), config.Port, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(common.Green("Running Client..."))
	fmt.Printf("%s | %s %s\n", common.Green("TCP"), common.Green("Remote Address:"), c.TCPConn.RemoteAddr().String())

	client.RunClient(c)
}

func serverModeInit() {
	config, err := getConfig("server_config.json")
	if err != nil {
		panic(err)
	}

	s, err := server.NewServer(net.ParseIP(config.Ip), config.Port)
	if err != nil {
		panic(err)
	}

	fmt.Println(common.Green("Running Server..."))
	fmt.Printf("%s | %s %s\n", common.Green("TCP"), common.Green("Listen On:"), s.TCPListener.Addr().String())
	server.RunServer(s)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `
Version: 1.0
Usage: socketTest [-h] [-m mode]
Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
