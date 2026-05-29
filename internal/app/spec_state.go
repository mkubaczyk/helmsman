package app

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

type StatePath struct {
	Path string `json:"path"`
}

type StateFiles struct {
	StateFiles []StatePath `json:"stateFiles"`
}

// specFromYAML reads a yaml file and decodes it to a state type.
// parser which throws an error if the YAML file is not valid.
func (pc *StateFiles) specFromYAML(file string) error {
	rawYamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("specFromYaml %v %v", file, err)
		return err
	}

	yamlFile := string(rawYamlFile)

	return yaml.Unmarshal([]byte(yamlFile), pc)
}
