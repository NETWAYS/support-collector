package base

import (
	"testing"
)

func TestGetKernelInfo(t *testing.T) {
	k, err := GetKernelInfo()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if k.Kernel == "" {
		t.Error("Expected Kernel to be not empty")
	}

	if k.Version == "" {
		t.Error("Expected Version to be not empty")
	}

	if k.Hostname == "" {
		t.Error("Expected Hostname to be not empty")
	}

	if k.FQDN == "" {
		t.Error("Expected FQDN to be not empty")
	}
}
