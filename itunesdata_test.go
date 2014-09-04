package fwutils

import (
	"testing"
	"strings"
)

func TestGetSoftwareURLFor(t *testing.T) {
	identifiers := [...]string{"iPhone4,1", "iPad3,1", "iPhone2,1", "iPhone5,1", "AppleTV2,1", "DeviceNotFound"}
	vm := NewiTunesVersionMaster()

	for _, identifier := range identifiers {

		t.Log("Testing identifier:", identifier)

		url, err := vm.GetSoftwareURLFor(identifier)

		if !strings.Contains(url, identifier) {
			if err == nil {
				t.Log("Identifier not present in software URL")
				t.Fail()
			}
		}		
	}
}