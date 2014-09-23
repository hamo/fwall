package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type serverConfig struct {
	ListenAddr string `json:"listen_addr"`
	ListenPort int    `json:"listen_port"`

	MasterKey     string `json:"master_key"`
	EncryptMethod string `json:"encrypt_method"`
	Tunnel        string `json:"tunnel"`

	// FIXME: a better name
	UserDB string `json:"userDB"`
}

func parseConfigFile(path string) (*serverConfig, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return doParseConfigFile(content)
}

func doParseConfigFile(content []byte) (*serverConfig, error) {
	sc := new(serverConfig)
	err := json.Unmarshal(content, sc)
	if err != nil {
		return nil, err
	}

	if err := sanityCheckConfigFile(sc); err != nil {
		return nil, err
	}

	return sc, nil
}

func sanityCheckConfigFile(sc *serverConfig) error {
	if sc == nil {
		panic("error")
	}

	// listen address: if empty, set to 0.0.0.0
	if sc.ListenAddr == "" {
		sc.ListenAddr = "0.0.0.0"
	}

	// listen port is required
	if sc.ListenPort == 0 {
		return fmt.Errorf("Listen port is required")
	}

	// Masterkey is required
	if sc.MasterKey == "" {
		return fmt.Errorf("MasterKey is required")
	}

	// encryption method is required
	if sc.EncryptMethod == "" {
		return fmt.Errorf("Encryption Method is required")
	}

	// tunnel is required
	if sc.Tunnel == "" {
		return fmt.Errorf("Tunnel is required")
	}

	// UserDB is required
	if sc.UserDB == "" {
		return fmt.Errorf("UserDB is required")
	}

	return nil
}
