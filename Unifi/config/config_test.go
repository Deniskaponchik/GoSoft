package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	testSlice := []string{"TEST", "PROD"}

	for _, mode := range testSlice {
		t.Logf(mode)
		actual, err := NewConfig(mode)
		if err != nil {
			t.Errorf("Incorrect result. %s", err)
		} else {
			if actual.BpmUrl != "" {
				t.Logf("BpmUrl: %s", actual.BpmUrl)
			}
			if actual.SoapUrl != "" {
				t.Logf("SoapUrl: %s", actual.SoapUrl)
			}
			if actual.GlpiITsupport != "" {
				t.Logf("GlpiITsupport: %s", actual.GlpiITsupport)
			}
		}
		t.Logf("")
	}

}
