// SPDX-License-Identifier: AGPL-3.0-only

package cfg_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/brainupdaters/drlm-agent/cfg"

	"github.com/brainupdaters/drlm-common/pkg/fs"
	logger "github.com/brainupdaters/drlm-common/pkg/log"
	"github.com/brainupdaters/drlm-common/pkg/test"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type TestCfgSuite struct {
	test.Test
}

func TestCfg(t *testing.T) {
	suite.Run(t, new(TestCfgSuite))
}

func (s *TestCfgSuite) AssertCfg() {
	home, err := homedir.Dir()
	s.Require().Nil(err)

	s.Equal(&cfg.DRLMAgentConfig{
		Core: cfg.DRLMAgentCoreConfig{
			Host:     "localhost",
			Port:     50051,
			TLS:      true,
			CertPath: "cert/server.crt",
		},
		Log: logger.Config{
			Level: "info",
			File:  filepath.Join(home, ".log/drlm/agent.log"),
		},
	}, cfg.Config)
}

func (s *TestCfgSuite) TestInit() {
	s.Run("should initialize the configuration correctly", func() {
		fs.FS = afero.NewMemMapFs()
		s.Nil(afero.WriteFile(fs.FS, "/etc/drlm/agent.toml", nil, 0644))

		cfg.Init("")

		s.AssertCfg()
	})

	s.Run("should initialize the configuration correctly with a specified config file", func() {
		fs.FS = afero.NewMemMapFs()
		s.Nil(afero.WriteFile(fs.FS, "/agent.toml", nil, 0644))

		cfg.Init("/agent.toml")

		s.AssertCfg()
	})

	s.Run("should reload the configuration correctly", func() {
		fs.FS = afero.NewOsFs()

		d, err := afero.TempDir(fs.FS, "", "drlm-agent-config-reload")
		s.Nil(err)
		defer fs.FS.RemoveAll(d)

		f := filepath.Join(d, "agent.toml")
		s.Nil(afero.WriteFile(fs.FS, f, nil, 0644))

		cfg.Init(f)

		s.AssertCfg()

		s.Nil(afero.WriteFile(fs.FS, f, []byte("[core]\nport = 8000"), 0644))

		// TODO: Change this to a non specified amount of time
		time.Sleep(1 * time.Second)

		s.Equal(8000, cfg.Config.Core.Port)
	})

	s.Run("should exit if there's an error reading the configuration file", func() {
		fs.FS = afero.NewMemMapFs()

		s.Exits(func() { cfg.Init("") })
	})

	s.Run("should exit if there's an error decoding the configuration", func() {
		fs.FS = afero.NewMemMapFs()
		s.Nil(afero.WriteFile(fs.FS, "/etc/drlm/agent.json", []byte("invalid config"), 0644))

		s.Exits(func() { cfg.Init("") })
	})
}
