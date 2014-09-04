package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type localConfig struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	MasterKey  string `json:"master_key"`
	Method     string `json:"method"`
	Tunnel     string `json:"tunnel"`

	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	Username     string `json:"username"`
	Password     string `json:"password"`

	WhiteList string `json:"whitelist"`
	BlackList string `json:"blacklist"`
}

func parseConfigFile(path string) (*localConfig, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return doParseConfigFile(content)
}

func doParseConfigFile(content []byte) (*localConfig, error) {
	lc := new(localConfig)
	err := json.Unmarshal(content, lc)
	if err != nil {
		return nil, err
	}

	if err := sanityCheckConfigFile(lc); err != nil {
		return nil, err
	}

	return lc, nil
}

func sanityCheckConfigFile(lc *localConfig) error {
	if lc == nil {
		panic("error")
	}

	// Server address and port are required
	if lc.Server == "" || lc.ServerPort == 0 {
		return fmt.Errorf("Server address and port are empty")
	}

	// Masterkey is required
	if lc.MasterKey == "" {
		return fmt.Errorf("MasterKey is required")
	}

	// encryption method is required
	if lc.Method == "" {
		return fmt.Errorf("Encryption Method is required")
	}

	if lc.Tunnel == "" {
		return fmt.Errorf("Tunnel is required")
	}

	if lc.LocalAddress == "" {
		lc.LocalAddress = "127.0.0.1"
	}

	if lc.LocalPort == 0 {
		lc.LocalPort = 1080
	}

	if lc.Username == "" || lc.Password == "" {
		return fmt.Errorf("Username and password are required")
	}

	return nil
}
