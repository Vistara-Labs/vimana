package config

import "github.com/BurntSushi/toml"

// type Spacecore struct {
// 	Modes    []string
// 	Binary   string
// 	Download string
// }

type Config struct {
	// Spacecores map[string]Spacecore
	Spacecores map[string]Spacecore `toml:"spacecores"`
}

type Spacecore map[string]Mode

type Mode struct {
	Binary   string `toml:"binary"`
	Download string `toml:"download"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	return config, err
}
