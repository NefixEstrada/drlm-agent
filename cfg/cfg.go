// SPDX-License-Identifier: AGPL-3.0-only

package cfg

import (
	"path/filepath"
	"strings"

	logger "github.com/brainupdaters/drlm-common/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// Config has the values of the configuration
var Config *DRLMAgentConfig

// DRLMAgentConfig is the configuration of an Agent of DRLM
type DRLMAgentConfig struct {
	Core  DRLMAgentCoreConfig  `mapstructure:"core"`
	Minio DRLMAgentMinioConfig `mapstructure:"minio"`
	Log   logger.Config        `mapstructure:"log"`
}

// DRLMAgentCoreConfig is the configuration related with the DRLM Core of an Agent of DRLM
type DRLMAgentCoreConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	TLS      bool   `mapstructure:"tls"`
	CertPath string `mapstructure:"cert_path"`
	Secret   string `mapstructure:"secret"`
}

// DRLMAgentMinioConfig is the configuration related wtih the Minio of a DRLM Agent
type DRLMAgentMinioConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	SSL       bool   `mapstructure:"ssl"`
	CertPath  string `mapstructure:"cert_path"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// v is the viper instance for the configuration
var v *viper.Viper

// Init prepares the configuration and reads it
func Init(fs afero.Fs, f string) {
	v = viper.New()
	v.SetFs(fs)
	SetDefaults()

	// If provided, use the configuration file
	if f != "" {
		v.SetConfigFile(f)
	}

	cfgNotFound := false
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading the configuration: %v", err)
		}

		log.Warnln("configuration file not found, using default and environment vaules")
		cfgNotFound = true
	}

	if err := v.Unmarshal(&Config); err != nil {
		log.Fatalf("error decoding the configuration: %v", err)
	}

	if cfgNotFound {
		return
	}

	v.WatchConfig()
	v.OnConfigChange(func(fsnotify.Event) {
		log.Info("reloading configuration...")

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("error reading the configuration: %v", err)
		}

		if err := v.Unmarshal(&Config); err != nil {
			log.Fatalf("error decoding the configuration: %v", err)
		}

		log.Info("configuration reloaded successfully")
	})

}

// SetDefaults sets the default configurations for Viper
func SetDefaults() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error getting the home directory: %v", err)
	}

	v.SetConfigName("agent")
	v.AddConfigPath(".")
	v.AddConfigPath(filepath.Join(home, ".drlm"))
	v.AddConfigPath(filepath.Join(home, ".config/drlm"))
	v.AddConfigPath("/etc/drlm")

	v.SetDefault("core", map[string]interface{}{
		"host":      "localhost",
		"port":      50051,
		"tls":       true,
		"cert_path": "cert/server.crt",
		"secret":    "secret",
	})
	v.SetDefault("minio", map[string]interface{}{
		"host":       "minio",
		"port":       9443,
		"ssl":        true,
		"cert_path":  "cert/minio.crt",
		"access_key": "drlm3minio",
		"secret_key": "drlm3minio",
	})
	v.SetDefault("log", map[string]interface{}{
		"level": "info",
		"file":  filepath.Join(home, ".log/drlm/agent.log"),
	})

	v.SetEnvPrefix("DRLM_AGENT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
