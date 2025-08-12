package stanza

import (
	"encoding/xml"

	"github.com/google/uuid"
	"github.com/gsvd/goeland-xmpp/address"
)

type PresenceType string
type PresenceShow string
type PresenceOption func(*Presence)

const (
	PresenceTypeAvailable    PresenceType = ""
	PresenceTypeUnavailable  PresenceType = "unavailable"
	PresenceTypeSubscribe    PresenceType = "subscribe"
	PresenceTypeSubscribed   PresenceType = "subscribed"
	PresenceTypeUnsubscribe  PresenceType = "unsubscribe"
	PresenceTypeUnsubscribed PresenceType = "unsubscribed"
	PresenceTypeProbe        PresenceType = "probe"
	PresenceTypeError        PresenceType = "error"
)

const (
	ShowAway PresenceShow = "away"
	ShowChat PresenceShow = "chat"
	ShowDND  PresenceShow = "dnd"
	ShowXA   PresenceShow = "xa"
)

type Presence struct {
	XMLName  xml.Name     `xml:"presence"`
	ID       string       `xml:"id,attr,omitempty"`
	Type     PresenceType `xml:"type,attr,omitempty"`
	From     string       `xml:"from,attr,omitempty"`
	To       string       `xml:"to,attr,omitempty"`
	Show     PresenceShow `xml:"show,omitempty"`
	Status   string       `xml:"status,omitempty"`
	Priority int          `xml:"priority,omitempty"`
}

func NewPresence(opts ...PresenceOption) *Presence {
	p := &Presence{
		ID:   uuid.NewString(),
		Type: PresenceTypeAvailable,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithPresenceType(t PresenceType) PresenceOption {
	return func(p *Presence) {
		p.Type = t
	}
}

func WithPresenceFrom(from address.Address) PresenceOption {
	return func(p *Presence) {
		p.From = from.String()
	}
}

func WithPresenceTo(to address.Address) PresenceOption {
	return func(p *Presence) {
		p.To = to.String()
	}
}

func WithPresenceShow(show PresenceShow) PresenceOption {
	return func(p *Presence) {
		p.Show = show
	}
}

func WithPresenceStatus(status string) PresenceOption {
	return func(p *Presence) {
		p.Status = status
	}
}

func WithPresencePriority(priority int) PresenceOption {
	return func(p *Presence) {
		p.Priority = priority
	}
}

func (t *PresenceType) UnmarshalXMLAttr(attr xml.Attr) error {
	validTypes := map[string]PresenceType{
		"available":    PresenceTypeAvailable,
		"unavailable":  PresenceTypeUnavailable,
		"subscribe":    PresenceTypeSubscribe,
		"subscribed":   PresenceTypeSubscribed,
		"unsubscribe":  PresenceTypeUnsubscribe,
		"unsubscribed": PresenceTypeUnsubscribed,
		"probe":        PresenceTypeProbe,
		"error":        PresenceTypeError,
	}

	if presenceType, valid := validTypes[attr.Value]; valid {
		*t = presenceType
	}

	return nil
}

func (t PresenceType) MarshalText() ([]byte, error) {
	switch t {
	case PresenceTypeAvailable,
		PresenceTypeUnavailable,
		PresenceTypeSubscribe,
		PresenceTypeSubscribed,
		PresenceTypeUnsubscribe,
		PresenceTypeUnsubscribed,
		PresenceTypeProbe,
		PresenceTypeError:
		return []byte(t), nil
	default:
		return []byte(nil), nil
	}
}
