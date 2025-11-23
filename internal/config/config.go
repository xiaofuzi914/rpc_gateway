package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig        `yaml:"server"`
	Chains map[string]ChainCfg `yaml:"chains"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type ChainCfg struct {
	RPCEndpoints []string `yaml:"rpc_endpoint"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if cfg.Server.Addr == "" {
		return nil, fmt.Errorf("server address is required")
	}

	if len(cfg.Chains) == 0 {
		return nil, fmt.Errorf("at least one chain configuration is required")
	}

	for name, c := range cfg.Chains {
		if len(c.RPCEndpoints) == 0 {
			return nil, fmt.Errorf("chain %s has no rpc_endpoints", name)
		}
	}

	return &cfg, nil
}
