package common

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// This is a secrets key to secret file mapping. We will attempt to read in from secret files before loading anything else.
var keyToSecretMapping = map[string]string{}
var keyToSecretMappingMutex = sync.Mutex{}

// Configs will populate viper with specified configs.
func Configs(secretLocations []string) error {
	if len(secretLocations) > 0 {
		secrets := GetAllSecrets()
		for key, secretFilename := range secrets {
			loadSecretFileIntoKey(key, secretFilename, secretLocations)
		}
	}

	return nil
}

// loadSecretFileIntoKey will attempt to load the contents of a secret file into the given key.
// If the secret file doesn't exist, we'll skip this.
func loadSecretFileIntoKey(key string, filename string, secretLocations []string) error {
	for _, secretLocation := range secretLocations {
		fullFilename := filepath.Join(secretLocation, filename)
		stat, err := os.Stat(fullFilename)
		if err == nil && !stat.IsDir() {
			data, err := ioutil.ReadFile(fullFilename)
			if err != nil {
				return fmt.Errorf("error loading secret file %s from location %s", filename, secretLocation)
			}
			viper.Set(key, strings.TrimSpace(string(data)))
			return nil
		}
	}

	return nil
}

// LoadConfigs loads secrets objects given the provided list of configs and a custom secrets
func LoadConfigs(secretLocationsString string) error {
	var secretLocations []string
	if secretLocationsString != "" {
		secretLocations = strings.Split(secretLocationsString, ",")
	}

	// Load configs
	if err := Configs(secretLocations); err != nil {
		return fmt.Errorf("error loading secrets: %v", err)
	}

	return nil
}

// RegisterSecret will register the secret filename that will be used for the corresponding Viper string.
func RegisterSecret(key string, secretFileName string) {
	keyToSecretMappingMutex.Lock()
	keyToSecretMapping[key] = secretFileName
	keyToSecretMappingMutex.Unlock()
}

// GetAllSecrets will return Viper secrets keys and their corresponding secret filenames.
func GetAllSecrets() map[string]string {
	return keyToSecretMapping
}
