package servercommon

import (
	"encoding/json"
	"io"
	"os"

	"github.com/maygli/dynamictags"
)

type ConfigurationReader struct {
	cfgData any
	dict    map[string]string
}

func NewConfigurationReader(fileName string, cfgDict map[string]string) (*ConfigurationReader, error) {
	cfg := ConfigurationReader{
		dict: cfgDict,
	}
	if fileName != "" {
		cfgFile, err := os.Open(fileName)
		if err == nil {
			defer cfgFile.Close()
			fileData, err := io.ReadAll(cfgFile)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(fileData, &cfg.cfgData)
			if err != nil {
				return nil, err
			}
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return &cfg, nil
}

func (cfgReader ConfigurationReader) ReadConfiguration(server interface{}, path string) error {
	proc, err := dynamictags.NewConfigurationProcessor(cfgReader.cfgData, path)
	if err != nil {
		return err
	}
	if cfgReader.dict != nil && len(cfgReader.dict) > 0 {
		proc.SetDictionary(cfgReader.dict)
	}
	err = proc.Process(server, nil)
	return err
}
