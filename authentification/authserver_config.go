package authentification

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/sethvargo/go-envconfig"
)

const (
	AUTH_CONFIG_FILE_NAME string = "authconfig.json"
)

func initServer(server *AuthServer) error {
	cfgFile, err := os.Open(AUTH_CONFIG_FILE_NAME)
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
	server.CookieAge = server.AcessTokenAge * 2
	return nil
}
