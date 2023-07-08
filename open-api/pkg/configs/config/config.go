package config

import (
	_ "github.com/spf13/viper"
	"sync"
)

type config struct {
	GRPCClientConfig GRPCClientConfig `mapstructure:"grpc_client"`
	Redis            Redis            `mapstructure:"redis"`
	Project          Project          `mapstructure:"project"`
}

type GRPCClientConfig struct {
	Clients map[string]GRPCClientInfo `mapstructure:"grpc_client"`
}

type GRPCClientInfo struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Db   string `mapstructure:"db"`
}
type Project struct {
	SecretPwd string `mapstructure:"secret_pwd"`
}

var (
	configs    *config
	configOnce sync.Once
)

// NewConfigs 单列入口
func NewConfigs() *config {
	configOnce.Do(func() {
		configs = &config{}
	})
	return configs
}
