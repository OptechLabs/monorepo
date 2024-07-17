package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadWoDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := LoadFromFile("config_test.json", false)
	assert.NoError(t, err)
	assert.Equal(t, "test", cfg.AppName)
	assert.Equal(t, "otoslocal.com", cfg.RootDomain)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.HTTPServerConfig.Port)
	assert.NotNil(t, cfg.GRPCClientConfigs["orders"])
	assert.Equal(t, "orders", cfg.GRPCClientConfigs["orders"].Host)
	assert.Equal(t, "8081", cfg.GRPCClientConfigs["orders"].Port)
	assert.Equal(t, "test:url", cfg.DBConfigs["primary"].ConnectionURL)
}

func Test_LoadWithDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := LoadFromFile("config_test_no_env.json", true)
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.HTTPServerConfig.Port)
	assert.Equal(t, 30, cfg.HTTPServerConfig.ShutdownWait)
	assert.Equal(t, 15, cfg.HTTPServerConfig.WriteTimeout)
	assert.Equal(t, 15, cfg.HTTPServerConfig.ReadTimeout)
	assert.Equal(t, 60, cfg.HTTPServerConfig.IdleTimeout)
	assert.NotNil(t, cfg.GRPCClientConfigs["orders"])
	assert.Equal(t, "orders", cfg.GRPCClientConfigs["orders"].Host)
	assert.Equal(t, "8081", cfg.GRPCClientConfigs["orders"].Port)
}

func Test_LoadFromStringWithoutDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := LoadFromString(`{
		"appName": "test", 
		"environment": "development", 
		"httpServerConfig": {"port": "8080"}, 
		"grpcClientConfigs": {
		  "orders": {"host": "orders", "port": "8081"}
		}
	  }
`, true)
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.HTTPServerConfig.Port)
	assert.Equal(t, 30, cfg.HTTPServerConfig.ShutdownWait)
	assert.Equal(t, 15, cfg.HTTPServerConfig.WriteTimeout)
	assert.Equal(t, 15, cfg.HTTPServerConfig.ReadTimeout)
	assert.Equal(t, 60, cfg.HTTPServerConfig.IdleTimeout)
	assert.NotNil(t, cfg.GRPCClientConfigs["orders"])
	assert.Equal(t, "orders", cfg.GRPCClientConfigs["orders"].Host)
	assert.Equal(t, "8081", cfg.GRPCClientConfigs["orders"].Port)
}

func Test_LoadFromStringWDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := LoadFromString(`{
		"httpServerConfig": {"port": "8080"}, 
		"grpcClientConfigs": {"orders": {"host": "orders", "port": "8081"}}
	  }
`, true)
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.HTTPServerConfig.Port)
	assert.Equal(t, 30, cfg.HTTPServerConfig.ShutdownWait)
	assert.Equal(t, 15, cfg.HTTPServerConfig.WriteTimeout)
	assert.Equal(t, 15, cfg.HTTPServerConfig.ReadTimeout)
	assert.Equal(t, 60, cfg.HTTPServerConfig.IdleTimeout)
	assert.NotNil(t, cfg.GRPCClientConfigs["orders"])
	assert.Equal(t, "orders", cfg.GRPCClientConfigs["orders"].Host)
	assert.Equal(t, "8081", cfg.GRPCClientConfigs["orders"].Port)
}

func Test_DBConfigWDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := LoadFromString(`{
		"appName": "test",
		"rootDomain": "otoslocal.com",
		"environment": "development",
		"httpServerConfig": {
		  "port": "8080"
		},
		"grpcClientConfigs": {
		  "orders": {
			"host": "orders",
			"port": "8081",
			"googleIAMAudience": "audience"
		  }
		},
		"dbConfigs": {
		  "primary": {
			"connectionURL": "test:url"
		  }
		}
	  }
`, true)
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.HTTPServerConfig.Port)
	assert.Equal(t, 30, cfg.HTTPServerConfig.ShutdownWait)
	assert.Equal(t, 15, cfg.HTTPServerConfig.WriteTimeout)
	assert.Equal(t, 15, cfg.HTTPServerConfig.ReadTimeout)
	assert.Equal(t, 60, cfg.HTTPServerConfig.IdleTimeout)
	assert.NotNil(t, cfg.GRPCClientConfigs["orders"])
	assert.Equal(t, "orders", cfg.GRPCClientConfigs["orders"].Host)
	assert.Equal(t, "8081", cfg.GRPCClientConfigs["orders"].Port)
	assert.Equal(t, "test:url", cfg.DBConfigs["primary"].ConnectionURL)
	assert.Equal(t, 5, cfg.DBConfigs["primary"].MaxIdleConns)
	assert.Equal(t, 10, cfg.DBConfigs["primary"].MaxOpenConns)
}
