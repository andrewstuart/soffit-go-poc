package soffit

import "github.com/blang/semver"

// Payload is the tentative name for the main Soffit request body
type Payload struct {
	Request Request `json:"request"`
}

// Request is the go representation of the Soffit JSON request format v_1.
type Request struct {
	Mode        string              `json:"mode"`
	Namespace   string              `json:"namespace"`
	WindowState string              `json:"windowState"`
	Portal      PortalInfo          `json:"portal"`
	Preferences map[string][]string `json:"preferences"`
	User        UserDetails         `json:"user"`
}

// PortalInfo is the representation of the portal information sent by the uPortal server.
type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
}

// UserDetails is the representation of the user information sent by uPortal.
type UserDetails struct {
	Username string `json:"username"`
}
