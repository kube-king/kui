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

func YamlToString(data interface{}) string {
	marshal, err := yaml.Marshal(data)
	if err != nil {
		return ""
	}

	return string(marshal)
}

func YamlToFile(data interface{}, file string, mode os.FileMode) error {
	marshal, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, marshal, mode)
	if err != nil {
		return err
	}
	return nil
}
