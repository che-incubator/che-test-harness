package util


import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	"io"
	"os"
	"path/filepath"
)

//##############Inspiration to create debug file#########
//https://github.com/onsi/ginkgo/blob/master/reporters/junit_reporter.go#L70
//#################################
// DetailsReporter is a ginkgo reporter which dumps information regarding the tests.
type DetailsReporter struct {
	Writer io.Writer
}

// NewDetailsReporterFile returns a reporter which will create the file given and dump the specs
// to it as they complete.
func NewDetailsReporterFile(filename string) *DetailsReporter {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		_ = fmt.Errorf("%#v\n", err)
		panic(err)
	}
	f, err := os.Create(absPath)
	if err != nil {
		_ = fmt.Errorf("%#v\n", err)
		panic(err)
	}
	return NewDetailsReporterWithWriter(f)
}

// NewDetailsReporterWithWriter returns a reporter which will write the SpecSummary objects as tests
// complete to the given writer.
func NewDetailsReporterWithWriter(w io.Writer) *DetailsReporter {
		return &DetailsReporter{
			Writer: w,
		}
}

// SpecSuiteWillBegin is implemented as a noop to satisfy the reporter interface for ginkgo.
func (reporter *DetailsReporter) SpecSuiteWillBegin(cfg config.GinkgoConfigType, summary *types.SuiteSummary) {
}

// SpecSuiteDidEnd is implemented as a noop to satisfy the reporter interface for ginkgo.
func (reporter *DetailsReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {}

// SpecDidComplete is invoked by Ginkgo each time a spec is completed (including skipped specs).
func (reporter *DetailsReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	b, err := json.Marshal(specSummary)
	if err != nil {
		_ = fmt.Errorf("Error in detail reporter: %v", err)
		return
	}
	_, err = reporter.Writer.Write(b)
	if err != nil {
		_ = fmt.Errorf("Error saving test details in detail reporter: %v", err)
		return
	}
	// Printing newline between records for easier viewing in various tools.
	_, err = fmt.Fprintln(reporter.Writer, "")
	if err != nil {
		_ = fmt.Errorf("Error saving test details in detail reporter: %v", err)
		return
	}
}

// SpecWillRun is implemented as a noop to satisfy the reporter interface for ginkgo.
func (reporter *DetailsReporter) SpecWillRun(specSummary *types.SpecSummary) {}

// BeforeSuiteDidRun is implemented as a noop to satisfy the reporter interface for ginkgo.
func (reporter *DetailsReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {}

// AfterSuiteDidRun is implemented as a noop to satisfy the reporter interface for ginkgo.
func (reporter *DetailsReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {
	if c, ok := reporter.Writer.(io.Closer); ok {
		_ = c.Close()
	}
}


