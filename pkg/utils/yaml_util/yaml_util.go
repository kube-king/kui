package yaml_util

import (
	"gopkg.in/yaml.v3"
	"os"
)

func UnYamlFile(file string, data interface{}) error {

	mHostYamlData, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(mHostYamlData, data)
	if err != nil {
		return err
	}

	return nil
}
