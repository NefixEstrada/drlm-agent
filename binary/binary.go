package binary

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/brainupdaters/drlm-agent/cfg"

	"github.com/brainupdaters/drlm-common/pkg/minio"
	sdk "github.com/minio/minio-go/v6"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

// Install installls a binary
func Install(bucket, name string) error {
	h, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("error getting the home directory: %v", err)
	}

	fs := afero.NewOsFs()

	f, err := fs.OpenFile(filepath.Join(h, ".bin", name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("error opening the binary: %v", err)
	}
	defer f.Close()

	minio, err := minio.NewSDK(fs, cfg.Config.Minio.Host, cfg.Config.Minio.Port, cfg.Config.Minio.AccessKey, cfg.Config.Minio.SecretKey, cfg.Config.Minio.SSL, cfg.Config.Minio.CertPath)
	if err != nil {
		return fmt.Errorf("error creating the minio client: %v", err)
	}

	o, err := minio.GetObject(bucket, name, sdk.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("error getting the binary: %v", err)
	}
	defer o.Close()

	if _, err := io.Copy(f, o); err != nil {
		return fmt.Errorf("error writting the binary: %v", err)
	}

	return nil
}
