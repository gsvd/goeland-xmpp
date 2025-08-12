package stanza

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	t.Run("message types", testMessageTypes)
	t.Run("marshal", testMessageMarshal)
	t.Run("unmarshal", testMessageUnmarshal)
}

func testMessageTypes(t *testing.T) {
	tests := []struct {
		name     string
		msgType  MessageType
		expected string
	}{
		{"chat", ChatMessage, "chat"},
		{"groupchat", GroupChatMessage, "groupchat"},
		{"headline", HeadlineMessage, "headline"},
		{"normal", NormalMessage, "normal"},
		{"error", ErrorMessage, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.msgType))
		})
	}
}

func testMessageMarshal(t *testing.T) {
	runUUID := uuid.NewString()

	tests := []struct {
		name     string
		message  Message
		expected string
	}{
		{
			name: "basic chat message with resource",
			message: Message{
				Lang:   "en",
				ID:     fmt.Sprintf("%s-msg0", runUUID),
				Type:   ChatMessage,
				From:   "user@example.com",
				To:     "friend@example.com/tablet.iOS-18_6",
				Body:   "Hello, world!",
				Thread: "thread1",
			},
			expected: `<message xml:lang="en" id="` + runUUID + `-msg0" type="chat" from="user@example.com" to="friend@example.com/tablet.iOS-18_6"><body>Hello, world!</body><thread>thread1</thread></message>`,
		},
		{
			name: "normal message without thread",
			message: Message{
				ID:   fmt.Sprintf("%s-msg1", runUUID),
				Type: NormalMessage,
				From: "sender@example.com",
				To:   "recipient@example.com",
				Body: "Test message",
			},
			expected: `<message id="` + runUUID + `-msg1" type="normal" from="sender@example.com" to="recipient@example.com"><body>Test message</body></message>`,
		},
		{
			name: "message without type should omit type attr",
			message: Message{
				ID:   fmt.Sprintf("%s-msg2", runUUID),
				From: "sender@example.com",
				To:   "recipient@example.com",
				Body: "No type",
			},
			expected: `<message id="` + runUUID + `-msg2" from="sender@example.com" to="recipient@example.com"><body>No type</body></message>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.message)
			require.NoError(t, err, "failed to marshal message")
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func testMessageUnmarshal(t *testing.T) {
	runUUID := uuid.NewString()

	tests := []struct {
		name     string
		xml      string
		expected Message
	}{
		{
			name: "basic chat message",
			xml:  `<message type="chat" id="` + runUUID + `-msg0" from="user@example.com" to="friend@example.com"><body>Hello!</body><thread>t1</thread></message>`,
			expected: Message{
				Type:   ChatMessage,
				ID:     fmt.Sprintf("%s-msg0", runUUID),
				From:   "user@example.com",
				To:     "friend@example.com",
				Body:   "Hello!",
				Thread: "t1",
			},
		},
		{
			name: "message without type",
			xml:  `<message id="` + runUUID + `-msg1" from="sender@example.com"><body>Test</body></message>`,
			expected: Message{
				ID:   fmt.Sprintf("%s-msg1", runUUID),
				From: "sender@example.com",
				Body: "Test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg Message
			err := xml.Unmarshal([]byte(tt.xml), &msg)
			require.NoError(t, err, "failed to unmarshal message")

			assert.Equal(t, tt.expected.Type, msg.Type, "type mismatch")
			assert.Equal(t, tt.expected.ID, msg.ID, "id mismatch")
			assert.Equal(t, tt.expected.From, msg.From, "from mismatch")
			assert.Equal(t, tt.expected.To, msg.To, "to mismatch")
			assert.Equal(t, tt.expected.Body, msg.Body, "body mismatch")
			assert.Equal(t, tt.expected.Thread, msg.Thread, "thread mismatch")
		})
	}
}
