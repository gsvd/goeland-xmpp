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

	// v0.1.0 - Supporting Bind only
	Bind *Bind `xml:"urn:ietf:params:xml:ns:xmpp-bind bind,omitempty"`
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

	NSRoster     = "jabber:iq:roster"
	NSDiscoInfo  = "http://jabber.org/protocol/disco#info"
	NSDiscoItems = "http://jabber.org/protocol/disco#items"
	NSPing       = "urn:xmpp:ping"
	NSVCard      = "vcard-temp"
	NSBind       = "urn:ietf:params:xml:ns:xmpp-bind"
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
		ID:   id.New(),
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

func WithBindResource(resource string) IQOption {
	return func(iq *IQ) {
		if iq.Bind == nil {
			iq.Bind = &Bind{}
		}
		iq.Bind.Resource = resource
	}
}

func WithBindAddress(a address.Address) IQOption {
	return func(iq *IQ) {
		if iq.Bind == nil {
			iq.Bind = &Bind{}
		}
		iq.Bind.Address = a.String()
	}
}

func WithBindAddressStr(a string) IQOption {
	return func(iq *IQ) {
		if iq.Bind == nil {
			iq.Bind = &Bind{}
		}
		iq.Bind.Address = a
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
