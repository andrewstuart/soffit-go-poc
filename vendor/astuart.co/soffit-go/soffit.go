// Package soffit provides data structures for the data types one can expect
// from the uPortal Soffit portlet. This should help enable more rapid
// iteration when working in a soffit environment.
package soffit

import (
	"bytes"
	"log"
	"net/url"

	"github.com/blang/semver"
)

// Payload is the tentative name for the main Soffit request body
type Payload struct {
	Request Request     `json:"request"`
	User    UserDetails `json:"user"`
	Context Context     `json:"context"`
}

// Request is the go representation of the Soffit JSON request format v_1.
type Request struct {
	Mode        string `json:"mode"`
	WindowID    string `json:"windowId"`
	Namespace   string `json:"namespace"`
	WindowState string `json:"windowState"`

	Properties  map[string]string `json:"properties"`
	Preferences url.Values        `json:"preferences"`
	Attributes  url.Values        `json:"attributes"`
}

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

// UserDetails is the representation of the user information sent by uPortal.
type UserDetails struct {
	Username   string     `json:"username"`
	Attributes url.Values `json:"attributes"`
	Roles      []string   `json:"roles"`
	Groups     []string   `json:"groups"`
}

// Context represents information about the portal creating the request
type Context struct {
	PortalInfo            PortalInfo `json:"portalInfo"`
	SupportedWindowStates []string   `json:"supportedWindowStates"`
	Attributes            url.Values `json:"attributes"`
}
