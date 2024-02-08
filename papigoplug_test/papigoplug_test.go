package papigoplug_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	plug "github.com/Dewberry/papigoplug/papigoplug"
	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    shutdown()
    os.Exit(code)
}

func TestLog(t *testing.T) {
    plug.Log.Info("Hello, World! From, papigoplug.")
}

func TestGoodParamsWithoutOptional(t *testing.T) {
    pseudoOSArgs := []string{`test`, `{"last": "Torvalds", "first": "Linus"}`}
    paramsExpect := map[string]interface{}{"last": "Torvalds", "first": "Linus"}
    t.Logf("parsing pseudoOSArgs=%s", pseudoOSArgs)
	params, err := plug.ParseInput(pseudoOSArgs, allowedParams)
	if err != nil {
        plug.Log.Error(err)
        t.Fatal(err)
	}
    t.Logf("pseudoOSArgs parsed into params=%s", params)
    if !cmp.Equal(params, paramsExpect) {
        err = fmt.Errorf("params != paramsExpect: %s != %s", params, paramsExpect)
        t.Fatal(err)
    }
}

func TestGoodParamsWithOptional(t *testing.T) {
    pseudoOSArgs := []string{`test`, `{"last": "Torvalds", "first": "Linus", "middle": "Benedict"}`}
    paramsExpect := map[string]interface{}{"last": "Torvalds", "first": "Linus", "middle": "Benedict"}
    t.Logf("parsing pseudoOSArgs=%s", pseudoOSArgs)
	params, err := plug.ParseInput(pseudoOSArgs, allowedParams)
	if err != nil {
        t.Fatal(err)
	}
    t.Logf("pseudoOSArgs parsed into params=%s", params)
    if !cmp.Equal(params, paramsExpect) {
        err = fmt.Errorf("params != paramsExpect: %s != %s", params, paramsExpect)
        t.Fatal(err)
    }
}

func TestPrintResults(t *testing.T) {
    testResults := map[string]interface{}{"foo": "bar"}
    var err error = plug.PrintResults(testResults)
    if err != nil {
        t.Fatal(err)
    }
}

func TestBadParamsMissingRequiredKey(t *testing.T) {
    pseudoOSArgs := []string{`test`, `{"last": "Torvalds"}`}
    t.Logf("parsing pseudoOSArgs=%s", pseudoOSArgs)
	params, err := plug.ParseInput(pseudoOSArgs, allowedParams)
	if err == nil {
        err = errors.Join(err, fmt.Errorf("pseudoOSArgs should have failed to parse, but err was nil. They parsed to params=%s", params))
        t.Fatal(err)
	}
}

func TestBadParamsUnexpectedKey(t *testing.T) {
    pseudoOSArgs := []string{`test`, `{"last": "Torvalds", "first": "Linux", "typo": true}`}
    t.Logf("parsing pseudoOSArgs=%s", pseudoOSArgs)
	params, err := plug.ParseInput(pseudoOSArgs, allowedParams)
	if err == nil {
        err = errors.Join(err, fmt.Errorf("pseudoOSArgs should have failed to parse, but err was nil. They parsed to params=%s", params))
        t.Fatal(err)
	}
}

// TODO: func TestS3(t *testing.T)


// e.g. these are both allowed:
// {"last": "Torvalds", "first": "Linus"}
// {"last": "Torvalds", "first": "Linus", "middle": "Benedict"}
var allowedParams = plug.PluginParams{
    Required: []string{"last", "first"},
    Optional: []string{"middle"},
}

func initTestLog (level string) {
    plug.InitLog(level)
}

// setup is called by TestMain at the beginning of every Test
func setup() {
    initTestLog("info")
}

// shutdown is called by TestMain at the end of every Test
func shutdown() {

}
