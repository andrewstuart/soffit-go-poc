package soffit

import "github.com/blang/semver"

type SoffitRequest struct {
	Mode        string              `json:"mode"`
	Namespace   string              `json:"namespace"`
	WindowState string              `json:"windowState"`
	Portal      PortalInfo          `json:"portal"`
	Preferences map[string][]string `json:"preferences"`
}

type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
}
