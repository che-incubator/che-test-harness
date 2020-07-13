package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type TestHarnessConfiguration struct {
	Version               string                `yaml:"version"`
	OSD                   bool                  `yaml:"osd"`
	Flavor                string                `yaml:"flavor"`
	KubernetesImagePuller KubernetesImagePuller `yaml:"ImagePullerNamespace"`
}

type KubernetesImagePuller struct {
	Namespace string `yaml:"Namespace"`
	Image     string `yaml:"Image"`
	PullerImages string `yaml:"PullerImages"`
}

var TestHarnessConfig = TestHarnessConfiguration{}

// TODO: Check if all values are setted in yaml
func ParseConfigurationFile() error {
	fileLocation, err := filepath.Abs("deploy/test-harness.yaml")
	if err != nil {
		fmt.Errorf("Failed to locate operator deployment yaml, %s", err)
		return err
	}
	yamlFile, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		fmt.Errorf("Failed to locate operator deployment yaml, %s", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &TestHarnessConfig)

	return err
}
