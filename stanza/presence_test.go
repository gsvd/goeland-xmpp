package stanza

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresence(t *testing.T) {
	t.Run("presence types", testPresenceTypes)
	t.Run("presence show states", testPresenceShowStates)
	t.Run("marshal", testPresenceMarshal)
	t.Run("unmarshal", testPresenceUnmarshal)
}

func testPresenceTypes(t *testing.T) {
	tests := []struct {
		name     string
		presType PresenceType
		expected string
	}{
		{"available", PresenceTypeAvailable, ""},
		{"unavailable", PresenceTypeUnavailable, "unavailable"},
		{"subscribe", PresenceTypeSubscribe, "subscribe"},
		{"subscribed", PresenceTypeSubscribed, "subscribed"},
		{"unsubscribe", PresenceTypeUnsubscribe, "unsubscribe"},
		{"unsubscribed", PresenceTypeUnsubscribed, "unsubscribed"},
		{"probe", PresenceTypeProbe, "probe"},
		{"error", PresenceTypeError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.presType))
		})
	}
}

func testPresenceShowStates(t *testing.T) {
	tests := []struct {
		name     string
		show     PresenceShow
		expected string
	}{
		{"away", ShowAway, "away"},
		{"chat", ShowChat, "chat"},
		{"dnd", ShowDND, "dnd"},
		{"xa", ShowXA, "xa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.show))
		})
	}
}

func testPresenceMarshal(t *testing.T) {
	runUUID := uuid.NewString()

	tests := []struct {
		name     string
		presence Presence
		expected string
	}{
		{
			name: "available presence with show and status no type attr",
			presence: Presence{
				Type:   PresenceTypeAvailable,
				ID:     fmt.Sprintf("%s-pres0", runUUID),
				From:   "user@example.com/desktop",
				Show:   ShowChat,
				Status: "Ready to chat",
			},
			expected: `<presence id="` + runUUID + `-pres0" from="user@example.com/desktop"><show>chat</show><status>Ready to chat</status></presence>`,
		},
		{
			name: "unavailable presence",
			presence: Presence{
				Type: PresenceTypeUnavailable,
				ID:   fmt.Sprintf("%s-pres1", runUUID),
				From: "user@example.com/mobile",
			},
			expected: `<presence id="` + runUUID + `-pres1" type="unavailable" from="user@example.com/mobile"></presence>`,
		},
		{
			name: "presence with priority",
			presence: Presence{
				Type:     PresenceTypeAvailable,
				ID:       fmt.Sprintf("%s-pres2", runUUID),
				From:     "user@example.com/tablet",
				Show:     ShowAway,
				Status:   "Away from keyboard",
				Priority: 5,
			},
			expected: `<presence id="` + runUUID + `-pres2" from="user@example.com/tablet"><show>away</show><status>Away from keyboard</status><priority>5</priority></presence>`,
		},
		{
			name: "subscribe presence",
			presence: Presence{
				Type: PresenceTypeSubscribe,
				ID:   fmt.Sprintf("%s-pres3", runUUID),
				From: "user@example.com",
				To:   "friend@example.com",
			},
			expected: `<presence id="` + runUUID + `-pres3" type="subscribe" from="user@example.com" to="friend@example.com"></presence>`,
		},
		{
			name: "presence with negative priority no type attr",
			presence: Presence{
				Type:     PresenceTypeAvailable,
				ID:       fmt.Sprintf("%s-pres4", runUUID),
				From:     "bot@example.com/automation",
				Status:   "Automated client - do not disturb",
				Priority: -1,
			},
			expected: `<presence id="` + runUUID + `-pres4" from="bot@example.com/automation"><status>Automated client - do not disturb</status><priority>-1</priority></presence>`,
		},
		{
			name: "presence with dnd show no type attr",
			presence: Presence{
				Type:   PresenceTypeAvailable,
				ID:     fmt.Sprintf("%s-pres5", runUUID),
				From:   "user@example.com/work",
				Show:   ShowDND,
				Status: "In a meeting",
			},
			expected: `<presence id="` + runUUID + `-pres5" from="user@example.com/work"><show>dnd</show><status>In a meeting</status></presence>`,
		},
		{
			name: "presence with xa show no type attr",
			presence: Presence{
				Type:   PresenceTypeAvailable,
				ID:     fmt.Sprintf("%s-pres6", runUUID),
				From:   "user@example.com/laptop",
				Show:   ShowXA,
				Status: "Extended away",
			},
			expected: `<presence id="` + runUUID + `-pres6" from="user@example.com/laptop"><show>xa</show><status>Extended away</status></presence>`,
		},
		{
			name: "presence with zero priority no type attr",
			presence: Presence{
				ID:       fmt.Sprintf("%s-pres7", runUUID),
				From:     "user@example.com/client",
				Priority: 0,
			},
			expected: `<presence id="` + runUUID + `-pres7" from="user@example.com/client"></presence>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.presence)
			require.NoError(t, err, "Failed to marshal presence")
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func testPresenceUnmarshal(t *testing.T) {
	runUUID := uuid.NewString()

	tests := []struct {
		name     string
		xml      string
		expected Presence
	}{
		{
			name: "available presence with show no type attr means available",
			xml:  `<presence id="` + runUUID + `-pres0" from="user@example.com/resource"><show>chat</show><status>Available</status></presence>`,
			expected: Presence{
				Type:   PresenceTypeAvailable, // Should be set to available when no type attr
				ID:     fmt.Sprintf("%s-pres0", runUUID),
				From:   "user@example.com/resource",
				Show:   ShowChat,
				Status: "Available",
			},
		},
		{
			name: "unavailable presence",
			xml:  `<presence type="unavailable" id="` + runUUID + `-pres1" from="user@example.com/mobile"></presence>`,
			expected: Presence{
				Type: PresenceTypeUnavailable,
				ID:   fmt.Sprintf("%s-pres1", runUUID),
				From: "user@example.com/mobile",
			},
		},
		{
			name: "presence with priority no type attr means available",
			xml:  `<presence id="` + runUUID + `-pres2" from="user@example.com"><show>away</show><priority>10</priority></presence>`,
			expected: Presence{
				Type:     PresenceTypeAvailable, // Should be set to available when no type attr
				ID:       fmt.Sprintf("%s-pres2", runUUID),
				From:     "user@example.com",
				Show:     ShowAway,
				Priority: 10,
			},
		},
		{
			name: "subscribe presence",
			xml:  `<presence type="subscribe" id="` + runUUID + `-pres3" from="requester@example.com" to="target@example.com"></presence>`,
			expected: Presence{
				Type: PresenceTypeSubscribe,
				ID:   fmt.Sprintf("%s-pres3", runUUID),
				From: "requester@example.com",
				To:   "target@example.com",
			},
		},
		{
			name: "presence with negative priority no type attr means available",
			xml:  `<presence id="` + runUUID + `-pres4"><priority>-5</priority><status>Bot</status></presence>`,
			expected: Presence{
				Type:     PresenceTypeAvailable, // Should be set to available when no type attr
				ID:       fmt.Sprintf("%s-pres4", runUUID),
				Status:   "Bot",
				Priority: -5,
			},
		},
		{
			name: "presence with dnd no type attr means available",
			xml:  `<presence id="` + runUUID + `-pres5"><show>dnd</show><status>Busy</status></presence>`,
			expected: Presence{
				Type:   PresenceTypeAvailable, // Should be set to available when no type attr
				ID:     fmt.Sprintf("%s-pres5", runUUID),
				Show:   ShowDND,
				Status: "Busy",
			},
		},
		{
			name: "presence with xa no type attr means available",
			xml:  `<presence id="` + runUUID + `-pres6"><show>xa</show></presence>`,
			expected: Presence{
				Type: PresenceTypeAvailable, // Should be set to available when no type attr
				ID:   fmt.Sprintf("%s-pres6", runUUID),
				Show: ShowXA,
			},
		},
		{
			name: "error presence type",
			xml:  `<presence type="error" id="` + runUUID + `-pres7" from="user@example.com" to="invalid@example.com"></presence>`,
			expected: Presence{
				Type: PresenceTypeError,
				ID:   fmt.Sprintf("%s-pres7", runUUID),
				From: "user@example.com",
				To:   "invalid@example.com",
			},
		},
		{
			name: "probe presence type",
			xml:  `<presence type="probe" id="` + runUUID + `-pres8" from="server@example.com" to="user@example.com"></presence>`,
			expected: Presence{
				Type: PresenceTypeProbe,
				ID:   fmt.Sprintf("%s-pres8", runUUID),
				From: "server@example.com",
				To:   "user@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pres Presence
			err := xml.Unmarshal([]byte(tt.xml), &pres)
			require.NoError(t, err, "Failed to unmarshal presence")

			assert.Equal(t, tt.expected.Type, pres.Type, "Type mismatch")
			assert.Equal(t, tt.expected.ID, pres.ID, "ID mismatch")
			assert.Equal(t, tt.expected.From, pres.From, "From mismatch")
			assert.Equal(t, tt.expected.To, pres.To, "To mismatch")
			assert.Equal(t, tt.expected.Show, pres.Show, "Show mismatch")
			assert.Equal(t, tt.expected.Status, pres.Status, "Status mismatch")
			assert.Equal(t, tt.expected.Priority, pres.Priority, "Priority mismatch")
		})
	}
}
