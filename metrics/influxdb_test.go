package metrics

import (
	"testing"
)

func TestRegisterInfluxDB(t *testing.T) {
	RegisterInfluxDB(nil)
	err := Write("s")
	if err != nil {
		t.Errorf(err.Error())
	}
}
