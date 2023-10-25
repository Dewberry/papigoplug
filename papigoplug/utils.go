package papigoplug

import (
	"encoding/json"
	"errors"
	"fmt"
)

// PluginParams must be instantiated by the program to define/constrain the allowed parameters.
type PluginParams struct {
	Required []string
	Optional []string
}

// ParseInput scans the provided arguments list (e.g. CLI args from os.Args). It expects that there is one real arg,
// i.e. that osArgs has length 2, that osArgs[0] is the file name and that osArgs[1] is a JSON dictionary. It marshals
// the JSON dictionary into a map and returns it after validating that the provided parameters satisfy the constraints
// defined by allowedParams (that no Required parameters are missing and that all other provided parameters are part of
// the Optional list).
func ParseInput(osArgs []string, allowedParams PluginParams) (providedParams map[string]interface{}, err error) {
	if len(osArgs) != 2 {
		err = fmt.Errorf("%d sys args provided (expected 2, e.g. 1 CLI arg)", len(osArgs))
		return
	}

	err = json.Unmarshal([]byte(osArgs[1]), &providedParams)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("Failed to parse provided arg as JSON key-val pairs: %q", osArgs[1]))
		return
	}

	var inputErrs []error

	var missingInputs []string
	for _, r := range allowedParams.Required {
		if !mapContainsString(providedParams, r) {
			missingInputs = append(missingInputs, r)
		}
	}
	if len(missingInputs) > 0 {
		inputErrs = append(inputErrs, fmt.Errorf("Missing required input keys: %q", missingInputs))
	}

	var unexpectedInputs []string
	for p := range providedParams {
		if !sliceContainsString(allowedParams.Required, p) && !sliceContainsString(allowedParams.Optional, p) {
			unexpectedInputs = append(unexpectedInputs, p)
		}
	}
	if len(unexpectedInputs) > 0 {
		inputErrs = append(inputErrs, fmt.Errorf("Unexpected input keys: %q", unexpectedInputs))
	}

	if len(inputErrs) > 0 {
		err = errors.Join(append(inputErrs, fmt.Errorf("Required: %q. Optional: %q. Provided: %q", allowedParams.Required, allowedParams.Optional, providedParams))...)
		return
	}

	return
}

// pluginResults is used to format the JSON string printed at the end of the program.
type pluginResults struct {
	Results map[string]interface{} `json:"plugin_results"`
}

// PrintResults encodes the provided map into a JSON string within key "plugin_results" and prints that.
// This function must be called at the end of the program.
func PrintResults(results map[string]interface{}) (err error) {
	bytes, err := json.Marshal(pluginResults{Results: results})
	if err != nil {
		return
	}
	fmt.Println(string(bytes))
	return
}

// mapContainsString returns true if s exists in the keys of the map, otherwise false.
func mapContainsString(dict map[string]interface{}, s string) bool {
	for key := range dict {
		if s == key {
			return true
		}
	}
	return false
}

// sliceContainsString returns true if s exists in the slice, otherwise false.
func sliceContainsString(slice []string, s string) bool {
	for _, val := range slice {
		if s == val {
			return true
		}
	}
	return false
}
