package stanza

import (
	"encoding/xml"

	"github.com/gsvd/goeland-xmpp/address"
	"github.com/gsvd/goeland-xmpp/internal/id"
)

type MessageType string
type MessageOption func(*Message)

type Message struct {
	XMLName xml.Name    `xml:"message"`
	Lang    string      `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	ID      string      `xml:"id,attr,omitempty"`
	Type    MessageType `xml:"type,attr,omitempty"`
	From    string      `xml:"from,attr,omitempty"`
	To      string      `xml:"to,attr,omitempty"`
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
		ID:   id.New(),
		Type: NormalMessage,
	}

	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

func WithMessageLang(lang string) MessageOption {
	return func(m *Message) {
		m.Lang = lang
	}
}

func WithMessageType(t MessageType) MessageOption {
	return func(m *Message) {
		m.Type = t
	}
}

func WithMessageFrom(from address.Address) MessageOption {
	return func(m *Message) {
		m.From = from.String()
	}
}

func WithMessageTo(to address.Address) MessageOption {
	return func(m *Message) {
		m.To = to.String()
	}
}

func WithMessageBody(body string) MessageOption {
	return func(m *Message) {
		m.Body = body
	}
}

func WithMessageThread(thread string) MessageOption {
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
	case ChatMessage,
		ErrorMessage,
		GroupChatMessage,
		HeadlineMessage,
		NormalMessage:
		return []byte(t), nil
	default:
		return []byte(NormalMessage), nil
	}
}
