// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/brainupdaters/drlm-agent/cfg"
	"github.com/brainupdaters/drlm-agent/cli"

	"github.com/brainupdaters/drlm-common/pkg/fs"
	logger "github.com/brainupdaters/drlm-common/pkg/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "drlm-agent",
	Short: "TODO",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		cli.Main()
	},
}

// Execute is the main function of the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", `configuration file to use instead of the defaults ("/etc/drlm/agent.toml", "~/.config/drlm/agent.toml", "~/.drlm/agent.toml", "./agent.toml")`)
}

func initConfig() {
	fs.Init()
	cfg.Init(cfgFile)
	logger.Init(cfg.Config.Log)
}
