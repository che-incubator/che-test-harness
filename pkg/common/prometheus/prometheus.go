package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	prometheusFileNamePattern string = "%s.%s.prometheus.prom"
	cheMetricName string = "che_metadata"
)

type Metrics struct {
	metricRegistry   *prometheus.Registry
	cheGatherer    *prometheus.GaugeVec
}

// NewMetrics creates a new prometheus object using the given secrets object.
func NewMetrics() *Metrics {
	// Set up Prometheus prometheus registry and gatherers
	metricRegistry := prometheus.NewRegistry()
	cheGatherer := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: cheMetricName,
		},
		[]string{"ci_provider", "metadata_name", "test_id", "job_id"},
	)
	metricRegistry.MustRegister(cheGatherer)

	return &Metrics{
		metricRegistry: metricRegistry,
		cheGatherer: cheGatherer,
	}
}

// WritePrometheusFile collects data and writes it out in the prometheus export file format (https://github.com/prometheus/docs/blob/master/content/docs/instrumenting/exposition_formats.md)
// Returns the prometheus file name.
func (m *Metrics) WritePrometheusFile(reportDir string) (string, error) {
	m.processJSONFile(m.cheGatherer, filepath.Join(reportDir, "addon-metadata.json"))

	prometheusFileName := fmt.Sprintf(prometheusFileNamePattern, "openshift-ci", "che-test-harness")
	output, err := m.registryToExpositionFormat()

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(reportDir, prometheusFileName), output, os.FileMode(0644))
	if err != nil {
		return "", err
	}

	return prometheusFileName, nil
}

// processJSONFile takes a JSON file and converts it into prometheus prometheus of the general format
func (m *Metrics) processJSONFile(gatherer *prometheus.GaugeVec, jsonFile string) (err error) {
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	var jsonOutput interface{}

	if err = json.Unmarshal(data, &jsonOutput); err != nil {
		return err
	}

	m.jsonToPrometheusOutput(gatherer, jsonOutput.(map[string]interface{}), []string{})

	return nil
}

// jsonToPrometheusOutput will take the JSON and write it into the gauge vector.
func (m *Metrics) jsonToPrometheusOutput(gatherer *prometheus.GaugeVec, jsonOutput map[string]interface{}, context []string) {
	for k, v := range jsonOutput {
		fullContext := append(context, k)
		switch jsonObject := v.(type) {
		case map[string]interface{}:
			m.jsonToPrometheusOutput(gatherer, jsonObject, fullContext)
		default:
			metadataName := strings.Join(fullContext, ".")
			stringValue := fmt.Sprintf("%v", jsonObject)
			// We're only concerned with tracking float values in Prometheus as they're the only thing we can measure
			if floatValue, err := strconv.ParseFloat(stringValue, 64); err == nil {
					gatherer.WithLabelValues("openshift-ci",
						metadataName,
						"che-test-harness",
						os.Getenv("BUILD_ID")).Add(floatValue)
			}
		}
	}
}

// registryToExpositionFormat takes all of the gathered prometheus and writes them out in the exposition format
func (m *Metrics) registryToExpositionFormat() ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := expfmt.NewEncoder(buf, expfmt.FmtText)
	metricFamilies, err := m.metricRegistry.Gather()

	if err != nil {
		return []byte{}, fmt.Errorf("error while gathering prometheus: %v", err)
	}

	for _, metricFamily := range metricFamilies {
		if err := encoder.Encode(metricFamily); err != nil {
			return []byte{}, fmt.Errorf("error encoding metric family: %v", err)
		}
	}

	return buf.Bytes(), nil
}
