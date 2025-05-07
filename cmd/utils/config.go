package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type InitConfig struct {
	EthAddress    string          `toml:"eth_address,omitempty"`
	EthPrivateKey string          `toml:"eth_private_key,omitempty"`
	Kvm           bool            `toml:"kvm"`
	CpuCount      int             `toml:"cpu_count"`
	RamSize       float64         `toml:"ram_size"`  // in GB
	DiskSize      float64         `toml:"disk_size"` // in GB
	InitDate      string          `toml:"init_date"`
	SpaceCore     string          `toml:"space_core"`
	Analytics     AnalyticsConfig `toml:"analytics"`
}

type AnalyticsConfig struct {
	Enabled bool `toml:"enabled"`
}

func LoadVimanaConfig(configFile string) (*InitConfig, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &InitConfig{
			Analytics: AnalyticsConfig{Enabled: true}, // defaults to true
		}, nil
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config InitConfig
	if err := toml.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *InitConfig, initPath string) error {
	// convert the configuration to TOML
	var buffer bytes.Buffer

	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(config)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(initPath), 0755)

	err = ioutil.WriteFile(initPath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
