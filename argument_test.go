package argument

import (
	"os"
	"reflect"
	"testing"
)

func Test_Instance(t *testing.T) {
	var argumentus = New()

	var expected = []string{"/bin/bash", "./mydir"}
	var callbackValues = make([][]string, 0)

	argumentus.Add("copy", "cp", func(values []string) {
		callbackValues = append(callbackValues, values)
	})

	os.Args = []string{"--copy", expected[0], expected[1], "-cp", "/overwrite/h", "./newdir"}
	argumentus.Start()

	if len(callbackValues) < 1 {
		t.Fatalf("callback not executed")
	}

	if len(callbackValues) > 1 {
		t.Fatalf("argument must not be overwritten")
	}

	if !reflect.DeepEqual(expected, callbackValues[0]) {
		t.Fatalf("expected %v | got %v", expected, callbackValues[0])
	}

	// empty os args.
	os.Args = nil
	argumentus.Start()
}

func Test_InstanceAdd(t *testing.T) {
	var argumentus = New()

	// no fullname err.
	var err = argumentus.Add("", "1", nil)
	if err != ErrNoFullName {
		t.Fatalf("expected ErrNoFullName, got: %v", err.Error())
	}

	// same names err.
	err = argumentus.Add("1", "1", nil)
	if err != ErrSameNames {
		t.Fatalf("expected ErrSameNames, got: %v", err.Error())
	}

	// no callback err.
	err = argumentus.Add("123", "1234", nil)
	if err != ErrNoCallback {
		t.Fatalf("expected ErrNoCallback, got: %v", err.Error())
	}
}
