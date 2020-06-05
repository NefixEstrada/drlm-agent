// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/brainupdaters/drlm-agent/cfg"
	"github.com/brainupdaters/drlm-agent/cli"

	logger "github.com/brainupdaters/drlm-common/pkg/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	join    bool
	fs      afero.Fs
	rootCmd = &cobra.Command{
		Use:   "drlm-agent",
		Short: "TODO",
		Long:  "TODO",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Main(fs, join)
		},
	}
)

// Execute is the main function of the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", `configuration file to use instead of the defaults ("/etc/drlm/agent.toml", "~/.config/drlm/agent.toml", "~/.drlm/agent.toml", "./agent.toml")`)
	rootCmd.PersistentFlags().BoolVarP(&join, "join", "j", false, `automatically try to join the DRLM Core. It has to be accepted from drlmctl`)
}

func initConfig() {
	fs = afero.NewOsFs()

	cfg.Init(fs, cfgFile)
	logger.Init(cfg.Config.Log)

}
