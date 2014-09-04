package fwutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"errors"
	"strings"
	"howett.net/plist"
)

type DeviceIdentifier string
type BuildNumber string

type IndividualBuild struct {
	BuildVersion     BuildNumber
	DocumentationURL string
	FirmwareURL      string
	FirmwareSHA1     string
	ProductVersion   string
}

type BuildInformation struct {
	Restore              *IndividualBuild
	Update               *IndividualBuild
	SameAs               BuildNumber
	OfferRestoreAsUpdate bool
}

type VersionWrapper struct {
	MobileDeviceSoftwareVersions map[DeviceIdentifier]map[BuildNumber]*BuildInformation
}

type iTunesVersionMaster struct {
	MobileDeviceSoftwareVersionsByVersion map[string]*VersionWrapper
}

// the URL for the XML
const iTunesVersionURL = "http://ax.phobos.apple.com.edgesuite.net/WebObjects/MZStore.woa/wa/com.apple.jingle.appserver.client.MZITunesClientCheck/version/"

func (vm *iTunesVersionMaster) GetSoftwareURLFor(device string) (url string, err error) {

	for _, deviceSoftwareVersions := range vm.MobileDeviceSoftwareVersionsByVersion {
		for identifier, builds := range deviceSoftwareVersions.MobileDeviceSoftwareVersions {
			if string(identifier) == device {
				for _, build := range builds {
					if build.Restore != nil {
						
						// don't return protected ones if we can avoid it
						if strings.Contains(build.Restore.FirmwareURL, "protected://") {
							continue
						}
						return build.Restore.FirmwareURL, nil
					}
				}
			}
		}
	}

	return "", errors.New("Unable to find identifier")
}

// creates a new iTunesVersionMaster struct, parsed and ready to use
func NewiTunesVersionMaster() *iTunesVersionMaster {
	resp, err := http.Get(iTunesVersionURL)
	if err != nil {
		panic(err)
	}

	document, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	vm := iTunesVersionMaster{}

	dec := plist.NewDecoder(bytes.NewReader(document))
	dec.Decode(&vm)

	return &vm
}


