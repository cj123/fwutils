// device is a library for getting information about
// apple's iOS devices
package fwutils

type Device struct {
	Identifier  string
	BDID        int
	BoardConfig string
	CPID        int
	Platform    string
	SCEP        int
	DeviceClass string
}

// Populates device given an identifier. Does so by finding a software URL for the identifier
func (d *Device) DeviceInfo(identifier string) (err error) {
	vm := NewiTunesVersionMaster()

	url, err := vm.GetSoftwareURLFor(identifier)
	if err != nil {
		
		return err
	}

	return d.DeviceInfoGivenURL(url)
}

// Populates device given the URL for the device
func (d *Device) DeviceInfoGivenURL(firmwareURL string) (err error) {
	ip := NewIPSW(firmwareURL)

	restore, err := ip.GetRestorePlist()

	// assume it's the first device, never seen a plist (yet) with more than one...
	device := restore.Devices[0]

	// there may be a better way than this, but I have yet to find it
	d.Identifier = restore.ProductType
	d.BDID       = device.BDID
	d.BoardConfig = device.BoardConfig
	d.CPID = device.CPID
	d.Platform = device.Platform
	d.SCEP = device.SCEP
	d.DeviceClass = restore.DeviceClass

	return err

}
