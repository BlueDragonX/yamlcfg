package yamlcfg

import (
	"syscall"
	"testing"
	"time"
)

func TestLookupSignal(t *testing.T) {
	testSuccess := func(name string, expect syscall.Signal) {
		if signal, ok := LookupSignal(name); ok {
			if signal != expect {
				t.Errorf("signal is invalid: %v != %v", signal, expect)
			}
		} else {
			t.Errorf("signal not found: %s", name)
		}
	}

	testFailure := func(name string) {
		if _, ok := LookupSignal(name); ok {
			t.Errorf("signal was found but shouldn't have been: %s", name)
		}
	}

	// Try a variety of valid signal names.
	testSuccess("SIGINT", syscall.SIGINT)
	testSuccess("\tSIGTERM      ", syscall.SIGTERM)
	testSuccess("Sigsegv", syscall.SIGSEGV)
	testSuccess("sigTrap  ", syscall.SIGTRAP)

	// Try a signal that doesn't exist.
	testFailure("SIGFART")
}

// Make a data map.
func makeData(key interface{}, value interface{}) map[interface{}]interface{} {
	data := make(map[interface{}]interface{}, 1)
	data[key] = value
	return data
}

func TestGetMapItem(t *testing.T) {
	key := "stuff"
	wantValue := "value of stuff"
	data := makeData(key, wantValue)
	if haveValue, ok := GetMapItem(data, key); ok {
		if haveValue != wantValue {
			t.Errorf("value \"%s\" is incorrect: \"%s\" != \"%s\"", key, haveValue, wantValue)
		}
	} else {
		t.Errorf("value \"%s\" was not found", key)
	}

	key = "missing"
	if _, ok := GetMapItem(data, key); ok {
		t.Errorf("value \"%s\" was found", key)
	}
}

func TestGetBool(t *testing.T) {
	key := "stuff"
	wantValue := true
	dflt := false
	data := makeData(key, wantValue)

	haveValue := GetBool(data, key, dflt)
	if haveValue != wantValue {
		t.Errorf("bool value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetBool(data, key, dflt)
	if haveValue != dflt {
		t.Errorf("bool value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetBoolAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid bool value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, "invalid")
	GetBool(data, key, true)
}

func TestGetString(t *testing.T) {
	key := "stuff"
	wantValue := "value"
	dflt := "default"
	data := makeData(key, wantValue)

	haveValue := GetString(data, key, dflt)
	if haveValue != wantValue {
		t.Errorf("string value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetString(data, key, dflt)
	if haveValue != dflt {
		t.Errorf("string value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetStringAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid string value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, true)
	GetString(data, key, "default")
}

func TestGetStringArray(t *testing.T) {
	key := "stuff"
	wantValue := []interface{}{"value"}
	dflt := []string{"default"}
	data := makeData(key, wantValue)

	haveValue := GetStringArray(data, key, dflt)
	if len(haveValue) != 1 && haveValue[0] != wantValue[0] {
		t.Errorf("string array value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetStringArray(data, key, dflt)
	if len(haveValue) != 1 && haveValue[0] != dflt[0] {
		t.Errorf("string array value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetStringArrayAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid string array value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, "invalid")
	GetStringArray(data, key, []string{"default"})
}

func TestGetInt(t *testing.T) {
	key := "stuff"
	wantValue := 12
	dflt := 36
	data := makeData(key, wantValue)

	haveValue := GetInt(data, key, dflt)
	if haveValue != wantValue {
		t.Errorf("int value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetInt(data, key, dflt)
	if haveValue != dflt {
		t.Errorf("int value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetIntAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid int value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, "invalid")
	GetInt(data, key, 42)
}

func TestGetDuration(t *testing.T) {
	key := "stuff"
	wantValue := 12
	dflt := 36 * time.Second
	data := makeData(key, wantValue)

	haveValue := GetDuration(data, key, dflt)
	if haveValue != time.Duration(wantValue)*time.Second {
		t.Errorf("duration value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetDuration(data, key, dflt)
	if haveValue != dflt {
		t.Errorf("duration value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetDurationAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid duration value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, "invalid")
	GetDuration(data, key, 42*time.Second)
}

func TestGetSignal(t *testing.T) {
	key := "stuff"
	dataValue := "SIGINT"
	wantValue := syscall.SIGINT
	dflt := syscall.SIGTERM
	data := makeData(key, dataValue)

	haveValue := GetSignal(data, key, dflt)
	if haveValue != wantValue {
		t.Errorf("signal value \"%s\" is incorrect: %v != %v", key, haveValue, wantValue)
	}

	key = "missing"
	haveValue = GetSignal(data, key, dflt)
	if haveValue != dflt {
		t.Errorf("signal value \"%s\" is not default: %v != %v", key, haveValue, dflt)
	}
}

func TestGetSignalAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid signal value did not assert")
		}
	}()

	key := "stuff"
	data := makeData(key, "invalid")
	GetSignal(data, key, 42)
}
