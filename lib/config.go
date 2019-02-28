package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/brainupdaters/drlm-common/logger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type DrlmagentConfig struct {
	Logging logger.LoggingConfig `mapstructure:"logging"`
}

var Config *DrlmagentConfig

func SetConfigDefaults() {
	logger.SetLoggingConfigDefaults("drlm-agent")
}

func InitConfig(c string) {
	if c != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".drlm-core" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".drlm-agent")
	}

	SetConfigDefaults()

	// Enable environment variables
	// ex.: DRLMAGENT_DRLMAPI_PORT=8000
	viper.SetEnvPrefix("DRLMAGENT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		panic("Unable to unmarshal config")
	}
}
