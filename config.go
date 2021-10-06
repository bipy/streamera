package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"streamera/common"
)

type configuration struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Type string `json:"type"`
}

func getConfig(configPath string) (*configuration, error) {
	jsonBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(common.Red("Load Config File Failed!"))
		return nil, err
	}

	config := configuration{}
	if err := json.Unmarshal(jsonBuf, &config); err != nil {
		fmt.Println(common.Red("Parse Json Failed!"))
		return nil, err
	}
	return &config, nil
}
