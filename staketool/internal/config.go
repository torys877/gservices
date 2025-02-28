package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is a full config from file
type Config struct {
	AppConfig AppConfig `yaml:"config"`
}

// AppConfig contains the main application settings.
type AppConfig struct {
	RPCURL          string `yaml:"rpc_url"`
	ChainID         string `yaml:"chain_id"`
	ContractAddress string `yaml:"contract_address"`
	StakeABI        string `yaml:"stake_abi"`
	WithdrawABI     string `yaml:"withdraw_abi"`
	ClaimABI        string `yaml:"claim_abi"`
	CheckRewardsABI string `yaml:"checkrewards_abi"`
}

// LoadConfig loads config yaml file in Config
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file '%s': %w", filename, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file '%s': %w", filename, err)
	}

	return &config, nil
}
