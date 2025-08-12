package stanza

import (
	"encoding/xml"

	"github.com/google/uuid"
	"github.com/gsvd/goeland-xmpp/address"
)

type MessageType string
type MessageOption func(*Message)

type Message struct {
	XMLName xml.Name    `xml:"message"`
	Type    MessageType `xml:"type,attr,omitempty"`
	ID      string      `xml:"id,attr,omitempty"`
	From    string      `xml:"from,attr,omitempty"`
	To      string      `xml:"to,attr,omitempty"`
	Lang    string      `xml:"xml:lang,attr,omitempty"`
	Body    string      `xml:"body,omitempty"`
	Thread  string      `xml:"thread,omitempty"`
}

const (
	ChatMessage      MessageType = "chat"
	ErrorMessage     MessageType = "error"
	GroupChatMessage MessageType = "groupchat"
	HeadlineMessage  MessageType = "headline"
	NormalMessage    MessageType = "normal"
)

func NewMessage(opts ...MessageOption) *Message {
	msg := &Message{
		Type: NormalMessage,
		ID:   uuid.NewString(),
	}

	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

func WithType(t MessageType) MessageOption {
	return func(m *Message) {
		m.Type = t
	}
}

func WithTo(to address.Address) MessageOption {
	return func(m *Message) {
		m.To = to.String()
	}
}

func WithFrom(from address.Address) MessageOption {
	return func(m *Message) {
		m.From = from.String()
	}
}

func WithBody(body string) MessageOption {
	return func(m *Message) {
		m.Body = body
	}
}

func WithThread(thread string) MessageOption {
	return func(m *Message) {
		m.Thread = thread
	}
}

func (t *MessageType) UnmarshalXMLAttr(attr xml.Attr) error {
	validTypes := map[string]MessageType{
		"normal":    NormalMessage,
		"chat":      ChatMessage,
		"error":     ErrorMessage,
		"groupchat": GroupChatMessage,
		"headline":  HeadlineMessage,
	}

	if msgType, valid := validTypes[attr.Value]; valid {
		*t = msgType
	} else {
		*t = NormalMessage
	}

	return nil
}

func (t MessageType) MarshalText() ([]byte, error) {
	switch t {
	case ChatMessage, ErrorMessage, GroupChatMessage, HeadlineMessage, NormalMessage:
		return []byte(t), nil
	default:
		return []byte(NormalMessage), nil
	}
}
