package servercommon

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/sethvargo/go-envconfig"
)

const (
	CONFIGURE_FOLDER = "config"
)

func ConfigureServer(server interface{}, fileName string) error {
	realFileName := path.Join(CONFIGURE_FOLDER, fileName)
	cfgFile, err := os.Open(realFileName)
	if err == nil {
		defer cfgFile.Close()
		fileData, err := io.ReadAll(cfgFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(fileData, server)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	err = envconfig.Process(context.Background(), server)
	if err != nil {
		return err
	}
	return nil
}
