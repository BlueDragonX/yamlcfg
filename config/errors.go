package config

import (
	"fmt"
)

type ParseError string

// NewParseError creates a new parse error with the default format.
func NewParseError(key string, value interface{}, what string) error {
	var msg string
	if key == "" {
		msg = fmt.Sprintf(`config failed to parse %+v as %s`, value, what)
	} else {
		msg = fmt.Sprintf(`config failed to parse {%s: %+v} as %s`, key, value, what)
	}
	return ParseError(msg)
}

func (err ParseError) Error() string {
	return string(err)
}

// AssertHasKeys panics with a ParseError if a key is missing from the provided YAML map.
func AssertHasKeys(data interface{}, keys []string, where string) {
	mapping := data.(map[interface{}](interface{}))
	for _, key := range keys {
		if _, ok := mapping[key]; !ok {
			panic(ParseError(fmt.Sprintf(`config key {%s: ?} not found in %s`, key, where)))
		}
	}
}

// AssertIsBool panics with a ParseError if value is not a boolean.
func AssertIsBool(key string, value interface{}) {
	if _, ok := value.(bool); !ok {
		panic(NewParseError(key, value, "bool"))
	}
}

// AssertIsString panics with a ParseError if value is not a string.
func AssertIsString(key string, value interface{}) {
	if _, ok := value.(string); !ok {
		panic(NewParseError(key, value, "string"))
	}
}

// AssertIsInt panics with a ParseError if value is not an int.
func AssertIsInt(key string, value interface{}) {
	switch value.(type) {
	case int:
	default:
		panic(NewParseError(key, value, "int"))
	}
}

// AssertIsArray panics with a ParseError if value is not an array.
func AssertIsArray(key string, value interface{}) {
	if _, ok := value.([]interface{}); !ok {
		panic(NewParseError(key, value, "array"))
	}
}

// AssertIsStringArray panics with a ParseError if value is not a string array.
func AssertIsStringArray(key string, value interface{}) {
	AssertIsArray(key, value)
	for _, item := range value.([]interface{}) {
		if _, ok := item.(string); !ok {
			panic(NewParseError(key, value, "string array"))
		}
	}
}

// AssertIsMap panics with a ParseError if value is not a map.
func AssertIsMap(key string, value interface{}) {
	if _, ok := value.(map[interface{}]interface{}); !ok {
		panic(NewParseError(key, value, "map"))
	}
}

// AssertIsStringMap panics with a ParseError if value is not a string/string map.
func AssertIsStringMap(key string, value interface{}) {
	AssertIsMap(key, value)
	for name, item := range value.(map[interface{}]interface{}) {
		if _, ok := name.(string); !ok {
			panic(NewParseError(key, value, "map key"))
		}
		if _, ok := item.(string); !ok {
			panic(NewParseError(key, value, "map value"))
		}
	}
}
