package main

import (
	"io/ioutil"
	"encoding/json"
)

type configFile struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	MasterKey  string `json:"master_key"`
	Method     string `json:"method"`

	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func parseConfigFile(path string) (*configFile, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return doParseConfigFile(content)
}

func doParseConfigFile(content []byte) (*configFile, error) {
	cf := new(configFile)
	err := json.Unmarshal(content, cf)
	if err != nil {
		return nil, err
	}

	return cf, nil
}
