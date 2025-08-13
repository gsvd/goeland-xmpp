package stanza

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/gsvd/goeland-xmpp/address"
	"github.com/gsvd/goeland-xmpp/internal/id"
)

type IQType string
type IQOption func(*IQ)

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	Lang    string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	ID      string   `xml:"id,attr"`
	Type    IQType   `xml:"type,attr"` // Required
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`

	// v0.1.0
	Bind    *Bind         `xml:"urn:ietf:params:xml:ns:xmpp-bind bind,omitempty"`
	Ping    *Ping         `xml:"urn:xmpp:ping ping,omitempty"`
	Version *VersionQuery `xml:"jabber:iq:version query,omitempty"`
}

type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

type VersionQuery struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name,omitempty"`
	Version string   `xml:"version,omitempty"`
	OS      string   `xml:"os,omitempty"`
}

type Bind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	Address  string   `xml:"jid,omitempty"`
}

const (
	IQTypeGet    IQType = "get"
	IQTypeSet    IQType = "set"
	IQTypeResult IQType = "result"
	IQTypeError  IQType = "error"

	NSPing    = "urn:xmpp:ping"
	NSVersion = "jabber:iq:version"
	NSBind    = "urn:ietf:params:xml:ns:xmpp-bind"
)

var (
	ErrInvalidIQType = errors.New("iq: invalid type")
)

func (t IQType) Valid() bool {
	switch t {
	case IQTypeGet, IQTypeSet, IQTypeResult, IQTypeError:
		return true
	default:
		return false
	}
}

// NewIQ creates a new IQ stanza.
// If the IQType is invalid, an error is returned.
func NewIQ(t IQType, opts ...IQOption) (*IQ, error) {
	if !t.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidIQType, t)
	}

	iq := &IQ{
		ID:   id.NewUUID(),
		Type: t,
	}

	for _, opt := range opts {
		opt(iq)
	}

	return iq, nil
}

func NewIQGet(opts ...IQOption) (*IQ, error)    { return NewIQ(IQTypeGet, opts...) }
func NewIQSet(opts ...IQOption) (*IQ, error)    { return NewIQ(IQTypeSet, opts...) }
func NewIQResult(opts ...IQOption) (*IQ, error) { return NewIQ(IQTypeResult, opts...) }
func NewIQError(opts ...IQOption) (*IQ, error)  { return NewIQ(IQTypeError, opts...) }

func WithIQLang(lang string) IQOption {
	return func(iq *IQ) {
		iq.Lang = lang
	}
}

func WithIQFrom(from string) IQOption {
	return func(iq *IQ) {
		iq.From = from
	}
}

func WithIQTo(to string) IQOption {
	return func(iq *IQ) {
		iq.To = to
	}
}

func WithPing() IQOption {
	return func(iq *IQ) {
		iq.Ping = &Ping{}
	}
}

func WithVersion(name, version, os string) IQOption {
	return func(iq *IQ) {
		iq.Version = &VersionQuery{
			Name:    name,
			Version: version,
			OS:      os,
		}
	}
}

func WithBind(resource string, address address.Address) IQOption {
	return func(iq *IQ) {
		iq.Bind = &Bind{
			Resource: resource,
			Address:  address.String(),
		}
	}
}

func (t *IQType) UnmarshalXMLAttr(attr xml.Attr) error {
	v := IQType(attr.Value)

	if !v.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidIQType, attr.Value)
	}

	*t = v

	return nil
}

func (t IQType) MarshalText() ([]byte, error) {
	return []byte(t), nil
}
