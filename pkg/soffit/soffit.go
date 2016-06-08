package soffit

import "github.com/blang/semver"

// Request is the go representation of the Soffit JSON request format v_1.
type Request struct {
	Mode        string              `json:"mode"`
	Namespace   string              `json:"namespace"`
	WindowState string              `json:"windowState"`
	Portal      PortalInfo          `json:"portal"`
	Preferences map[string][]string `json:"preferences"`
}

// PortalInfo is the representation of the portal information sent by the uPortal server.
type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
}
