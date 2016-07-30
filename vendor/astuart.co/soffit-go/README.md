# soffit
--
    import "astuart.co/soffit-go"

Package soffit provides data structures for the data types one can expect from
the uPortal Soffit portlet. This should help enable more rapid iteration when
working in a soffit environment.

## Usage

#### type Context

```go
type Context struct {
	PortalInfo            PortalInfo `json:"portalInfo"`
	SupportedWindowStates []string   `json:"supportedWindowStates"`
	Attributes            url.Values `json:"attributes"`
}
```

Context represents information about the portal creating the request

#### type Payload

```go
type Payload struct {
	Request Request     `json:"request"`
	User    UserDetails `json:"user"`
	Context Context     `json:"context"`
}
```

Payload is the tentative name for the main Soffit request body

#### type PortalInfo

```go
type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
	Snapshot bool           `json:"snapshot"`
}
```

PortalInfo is the representation of the portal information sent by the uPortal
server.

#### func (*PortalInfo) UnmarshalJSON

```go
func (p *PortalInfo) UnmarshalJSON(bs []byte) error
```
UnmarshalJSON implements json.Unmarshaler

#### type Request

```go
type Request struct {
	Mode        string `json:"mode"`
	WindowID    string `json:"windowId"`
	Namespace   string `json:"namespace"`
	WindowState string `json:"windowState"`

	Properties  map[string]string `json:"properties"`
	Preferences url.Values        `json:"preferences"`
	Attributes  url.Values        `json:"attributes"`
}
```

Request is the go representation of the Soffit JSON request format v_1.

#### type UserDetails

```go
type UserDetails struct {
	Username   string     `json:"username"`
	Attributes url.Values `json:"attributes"`
	Roles      []string   `json:"roles"`
	Groups     []string   `json:"groups"`
}
```

UserDetails is the representation of the user information sent by uPortal.
