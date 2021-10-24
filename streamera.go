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

func main() {
	if Config.Mode == "" || (Config.Mode != "server" && Config.Mode != "client") {
		usage()
		return
	}
	defer fmt.Println(common.Green("Bye"))
	if Config.Mode == "server" {
		serverModeInit()
	} else {
		clientModeInit()
	}
}

func clientModeInit() {
	c, err := client.NewClient(net.ParseIP(Config.Ip), Config.Port, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(common.Green("Running Client..."))
	fmt.Printf("%s | %s %s\n", common.Green("TCP"), common.Green("Remote Address:"), c.TCPConn.RemoteAddr().String())

	client.RunClient(c)
}

func serverModeInit() {
	s, err := server.NewServer(net.ParseIP(Config.Ip), Config.Port)
	if err != nil {
		panic(err)
	}

	fmt.Println(common.Green("Running Server..."))
	fmt.Printf("%s | %s %s\n", common.Green("TCP"), common.Green("Listen On:"), s.TCPListener.Addr().String())
	server.RunServer(s)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `
Version: 1.2
Usage: socketTest [-h] [-c config]
Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
