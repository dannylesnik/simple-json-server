package models

import (
	"testing"
	"time"
)

func TestGetIsAliveResponse(t *testing.T) {
	currentTime := time.Now().UnixNano()

	isaliveResponse := GetIsAliveResponse()

	if currentTime > isaliveResponse.Timestemp {
		t.Error("timestemp is incorrect")
	}

	if isaliveResponse.Hostname == "" {
		t.Error("Hostname is incorrect")
	}

	if isaliveResponse.IP == "" {
		t.Error("IP is incorrect")
	}

	if isaliveResponse.VER == "" {
		t.Error("Version is incorrect")
	}

}
