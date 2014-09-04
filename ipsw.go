// ipsw is for doing things with IPSWs
package fwutils

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/DHowett/ranger"
	"howett.net/plist"
	"net/url"
	"net/http"
	"io"
	"strconv"
	"strings"
	"time"
	//"io/ioutil"
	//"github.com/davecgh/go-spew/spew"
)

type IPSW struct {
	DownloadURL string
	Properties *IPSWProperties
}

type IndividualDevice struct {
	Identifier  string
	BDID        int
	BoardConfig string
	CPID        int
	Platform    string
	SCEP        int
}

type Restore struct {
	DeviceClass         string
	Devices             []*IndividualDevice `plist:"DeviceMap"`
	ProductBuildVersion string
	ProductType         string
	ProductVersion      string
}

func ReadFile(reader *zip.Reader, file string) (result []byte, err error) {

	for _, f := range reader.File {
		//fmt.Println(f)
		if f.Name == file {
			data := make([]byte, f.UncompressedSize64)
			rc, err := f.Open()

			if err != nil {
				return nil, err
			}

			io.ReadFull(rc, data)
			rc.Close()

			return data, nil
		}
	}

	return nil, errors.New("Unable to find file")
}

// Get the restore Plist from the IPSW
func (ip *IPSW) GetRestorePlist() (parsed Restore, err error) {
	url, err := url.Parse(ip.DownloadURL)

	// retrieve the Restore.plist from the IPSW, possibly need the BuildManifest too?
	reader, err := ranger.NewReader(&ranger.HTTPRanger{URL: url})
	zipreader, err := zip.NewReader(reader, reader.Length())

	restorePlist, err := ReadFile(zipreader, "Restore.plist")
	// decode the plist
	dec := plist.NewDecoder(bytes.NewReader(restorePlist))

	dec.Decode(&parsed)

	return parsed, err
}

type IPSWProperties struct {
	Identifier string 
	Version    string 
	BuildID    string 
	Filename   string 
	Size       int64 
	MD5sum     string 
	SHA1sum    string 
	UploadDate string
	ReleaseDate string
	AppleTVSoftwareVersion string

	// ... and probably some more that i've not thought of yet
}


// given a longForm to parse
const longForm = "Mon, 2 Jan 2006 15:04:05 MST"

// PopulateInfo gets the information about the IPSW from the path
func (ip *IPSW) PopulateInfo() (err error) {
	restorePlist, err := ip.GetRestorePlist()

	// get the headers from the URL
	res, err := http.Get(ip.DownloadURL)
	
	// get the size as an int
	size, err := strconv.ParseInt(res.Header["Content-Length"][0], 10, 64)

	// the md5sum is the string before the `:` after the `"` has been removed
	md5sum := strings.Split(strings.Replace(res.Header["Etag"][0], `"`, "", -1), ":")[0]


	uploadDate, _ := time.Parse(longForm, res.Header["Last-Modified"][0])


	ip.Properties = &IPSWProperties{
		Identifier: restorePlist.ProductType,
		Version: restorePlist.ProductVersion,
		BuildID: restorePlist.ProductBuildVersion,
		Filename: restorePlist.ProductType + "_" + restorePlist.ProductVersion + "_" + restorePlist.ProductBuildVersion + "_Restore.ipsw",
		Size: size,
		MD5sum: md5sum,
		UploadDate: uploadDate.Format(longForm),
	}

	return err
}

func (ip *IPSW) PopulateInfoGivenBuild(build *IndividualBuild) {
	ip.Properties.SHA1sum = build.FirmwareSHA1

	t := time.Now()
	ip.Properties.ReleaseDate =	t.Format(longForm)
	if ip.Properties.Version != build.ProductVersion {
		// this is most likely an AppleTV, which incorrectly reports versions
		ip.Properties.AppleTVSoftwareVersion = build.ProductVersion
	}

	return
}

// creates a new IPSW given an URL
func NewIPSW(url string) *IPSW {
	return &IPSW{DownloadURL: url}
}