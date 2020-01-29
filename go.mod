module github.com/brainupdaters/drlm-agent

// TODO: Remove this when https://github.com/kahing/goofys/pull/462 gets merged
replace github.com/kahing/goofys => github.com/nefixestrada/goofys v0.0.0-20200129082319-8f9a979f4315

replace github.com/jacobsa/fuse => github.com/kahing/fusego v0.0.0-20191210234239-374cf4208103

require (
	github.com/brainupdaters/drlm-common v0.0.0-20200127123945-18d8b3139fff
	github.com/fsnotify/fsnotify v1.4.7
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	google.golang.org/grpc v1.24.0
)

go 1.13
