package fwutils

import (

	"testing"
	//"github.com/davecgh/go-spew/spew"
)

func TestDeviceInfo(t *testing.T) {
	identifier := "iPhone3,1"

	d := Device{}
	d.DeviceInfo("iPhone4,1")
	
	t.Log(d)

	if d.Identifier != identifier {
		t.Log(d)
	}
}