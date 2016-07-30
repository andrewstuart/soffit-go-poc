// Package soffit provides data structures for the data types one can expect
// from the uPortal Soffit portlet. This should help enable more rapid
// iteration when working in a soffit environment.
package soffit

import "net/url"

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
