package parameters

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cnabio/cnab-go/bundle"
	"gopkg.in/yaml.v2"
)

// PorterInternal is a string that can be used to designate a parameter
// as internal to Porter
const PorterInternal = "porter-internal"

// ParseVariableAssignments converts a string array of variable assignments
// into a map of keys and values
// Example:
// [a=b c=abc1232=== d=banana d=pineapple] becomes map[a:b c:abc1232=== d:[pineapple]]
func ParseVariableAssignments(params []string) (map[string]string, error) {
	variables := make(map[string]string)
	for _, p := range params {

		parts := strings.SplitN(p, "=", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid parameter (%s), must be in name=value format", p)
		}

		variable := strings.TrimSpace(parts[0])
		if variable == "" {
			return nil, fmt.Errorf("invalid parameter (%s), variable name is required", p)
		}
		value := strings.TrimSpace(parts[1])

		variables[variable] = value
	}

	return variables, nil
}

// Load a ParameterSet from a file at a given path.
//
// It does not load the individual parameters.
func Load(path string) (ParameterSet, error) {
	var pset ParameterSet
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return pset, err
	}
	return pset, yaml.Unmarshal(data, &pset)
}

// IsInternal determines if the provided param is an internal parameter
// to Porter after analyzing the provided bundle
func IsInternal(param string, bun bundle.Bundle) bool {
	if param, exists := bun.Parameters[param]; exists {
		if def, exists := bun.Definitions[param.Definition]; exists {
			return def.Comment == PorterInternal
		}
	}
	return false
}
