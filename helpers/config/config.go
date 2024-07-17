package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	AppName           string                  `json:"appName" validate:"required"`
	RootDomain        string                  `json:"rootDomain" validate:"required"`
	SessionKey        string                  `json:"sessionKey"`
	Environment       string                  `json:"environment" validate:"required,oneof=development staging production"`
	GoogleProjectID   string                  `json:"googleProjectID"`
	HTTPServerConfig  ServerConfig            `json:"httpServerConfig"`
	GRPCServerConfig  ServerConfig            `json:"grpcServerConfig"`
	GRPCClientConfigs map[string]ClientConfig `json:"grpcClientConfigs"` //map[name]ClientConfig
	DBConfigs         map[string]DBConfig     `json:"dbConfigs"`         //map[use]DBConfig: ex. map["main"]DBConfig, map["readOnly"]DBConfig
	PubSubConfig      PubSubConfig            `json:"pubSubConfig"`
	AUTH0Config       Auth0Config             `json:"auth0Config"`
}

type Auth0Config struct {
	Domain       string `json:"domain"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type ClientConfig struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	Port              string `json:"port"`
	GoogleIAMAudience string `json:"googleIAMAudience"`
}

type ServerConfig struct {
	Port         string `json:"port"`
	ShutdownWait int    `json:"shutdownWait"`
	WriteTimeout int    `json:"writeTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	IdleTimeout  int    `json:"idleTimeout"`
}

type DBConfig struct {
	ConnectionURL string `json:"connectionURL"`
	MaxIdleConns  int    `json:"maxIdleConns"`
	MaxOpenConns  int    `json:"maxOpenConns"`
}

type PubSubConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadFromFile(configFile string, useDefaults bool) (config Config, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer f.Close()
	configBytes, err := io.ReadAll(f)
	if err != nil {
		return config, err
	}

	if useDefaults {
		return LoadConfigWithDefaults(configBytes)
	}
	return LoadConfig(configBytes)
}

func LoadFromString(configValues string, useDefaults bool) (config Config, err error) {
	if useDefaults {
		return LoadConfigWithDefaults([]byte(configValues))
	}
	return LoadConfig([]byte(configValues))
}

func LoadConfig(configValues []byte) (Config, error) {
	cfg, err := load(&Config{
		HTTPServerConfig: ServerConfig{},
		GRPCServerConfig: ServerConfig{},
	}, configValues)
	return *cfg, err
}

func LoadConfigWithDefaults(configValues []byte) (Config, error) {
	cfg, err := load(&Config{
		Environment: "development",
		HTTPServerConfig: ServerConfig{
			ShutdownWait: 30,
			WriteTimeout: 15,
			ReadTimeout:  15,
			IdleTimeout:  60,
		},
	}, configValues)

	if cfg.DBConfigs != nil {
		for connectionName, dbConfig := range cfg.DBConfigs {
			cfg.DBConfigs[connectionName] = DBConfig{
				ConnectionURL: dbConfig.ConnectionURL,
				MaxIdleConns: func(in int) int {
					if in == 0 {
						return 5
					}
					return in
				}(dbConfig.MaxIdleConns),
				MaxOpenConns: func(in int) int {
					if in == 0 {
						return 10
					}
					return in
				}(dbConfig.MaxOpenConns),
			}
		}
	}

	return *cfg, err
}

func load(cfg *Config, configValues []byte) (*Config, error) {
	if err := json.Unmarshal(configValues, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
