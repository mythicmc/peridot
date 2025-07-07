package config

type Config struct {
	Location string `json:"location"`
	// Supported: "vanilla", "paper", "velocity"
	Type    string   `json:"type"`
	Repos   []string `json:"repos"`
	Plugins []string `json:"plugins"`
}

type Configs map[string]Config
