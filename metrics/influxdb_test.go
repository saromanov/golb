package metrics

import (
	"testing"
)

func TestRegisterInfluxDB(t *testing.T) {
	RegisterInfluxDB()
	err := Write("s")
	if err != nil {
		t.Errorf(err.Error())
	}
}
