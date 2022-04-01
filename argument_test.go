package argument

import (
	"testing"
)

// add / error
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
