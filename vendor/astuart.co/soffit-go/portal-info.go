package soffit

import (
	"bytes"
	"log"

	"github.com/blang/semver"
)

// PortalInfo is the representation of the portal information sent by the uPortal server.
type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
	Snapshot bool           `json:"snapshot"`
}

var snapshotBytes = []byte("-SNAPSHOT")

// UnmarshalJSON implements json.Unmarshaler
func (p *PortalInfo) UnmarshalJSON(bs []byte) error {
	if len(bs) < 2 {
		return nil
	}

	// Reslice, removing quotations from ends
	if bs[0] == '"' {
		bs = bs[1:]
	}
	if bs[len(bs)-1] == '"' {
		bs = bs[:len(bs)-1]
	}

	bss := bytes.Split(bs, []byte{'/'})

	if len(bss) == 0 {
		return nil
	}

	p.Provider = string(bss[0])

	if len(bss) < 2 {
		return nil
	}

	sv := bss[1]

	if bytes.Contains(sv, snapshotBytes) {
		p.Snapshot = true
		sv = bytes.Replace(sv, snapshotBytes, []byte{}, -1)
	}

	log.Println(string(sv))

	var err error
	p.Version, err = semver.Parse(string(sv))
	return err
}
