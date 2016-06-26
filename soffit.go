package main

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/blang/semver"
)

// Payload is the tentative name for the main Soffit request body
type Payload struct {
	Request Request     `json:"request"`
	User    UserDetails `json:"user"`
}

// Request is the go representation of the Soffit JSON request format v_1.
type Request struct {
	Mode        string `json:"mode"`
	Namespace   string `json:"namespace"`
	WindowState string `json:"windowState"`

	Properties  map[string]string   `json:"properties"`
	Preferences map[string][]string `json:"preferences"`
	// Portal      PortalInfo          `json:"portal"`
}

// PortalInfo is the representation of the portal information sent by the uPortal server.
type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
}

// UserDetails is the representation of the user information sent by uPortal.
type UserDetails struct {
	Username   string              `json:"username"`
	Session    Session             `json:"session"`
	Attributes map[string][]string `json:"attributes"`
	Roles      []string            `json:"roles"`
}

// Session is the representation of the user session
type Session struct {
	CreationTime time.Time     `json:"creationTime"`
	TTL          time.Duration `json:"maxInactiveInterval"`
}

type sInter struct {
	CreationTime        int64
	MaxInactiveInterval int
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Session) UnmarshalJSON(bs []byte) error {
	var si sInter

	err := json.NewDecoder(bytes.NewReader(bs)).Decode(&si)
	if err != nil {
		return err
	}

	log.Println(si)

	s.CreationTime = time.Unix(si.CreationTime/1000, (si.CreationTime%1000)*1000000)
	s.TTL = time.Duration(si.MaxInactiveInterval) * time.Second

	return nil
}
