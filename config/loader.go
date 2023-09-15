package config

import "github.com/BurntSushi/toml"

// type Component struct {
// 	Modes    []string
// 	Binary   string
// 	Download string
// }

type Config struct {
	// Components map[string]Component
	Components map[string]Component `toml:"components"`
}

type Component map[string]Mode

type Mode struct {
	Binary   string `toml:"binary"`
	Download string `toml:"download"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	return config, err
}
