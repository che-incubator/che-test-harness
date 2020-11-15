package context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type metadata struct {
	Setup   `json:"setup"`
	Metrics `json:"metrics"`
}

var TestInstance = metadata{}

// WriteToJSON will marshall the metadata struct and write it into the given file.
func (m *metadata) WriteToJSON(outputFilename string) (err error) {
	var data []byte
	if data, err = json.Marshal(m); err != nil {
		return err
	}
	fmt.Println(m.CheNamespace)
	if err = ioutil.WriteFile(outputFilename, data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}
