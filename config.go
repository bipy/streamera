package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"streamera/common"
)

var Config *Configuration

type Configuration struct {
	Mode string `json:"mode"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func init() {
	flag.Usage = usage
	modePtr := flag.String("m", "", "run mode: \"server\" or \"client\"")
	hostPtr := flag.String("h", "127.0.0.1", "Host")
	portPtr := flag.Int("p", 6666, "Port")
	configPathPtr := flag.String("c", "", "Specific configuration file")
	flag.Parse()

	if *configPathPtr == "" {
		Config = &Configuration{
			Mode: *modePtr,
			Ip:   *hostPtr,
			Port: *portPtr,
		}
	} else {
		Config = parseJsonConfig(*configPathPtr)
	}
}

func parseJsonConfig(configPath string) *Configuration {
	jsonBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(common.Red("Load Config File Failed!"))
		panic(err)
	}

	c := Configuration{}
	if err := json.Unmarshal(jsonBuf, &c); err != nil {
		fmt.Println(common.Red("Parse Json Failed!"))
		panic(err)
	}
	return &c
}
