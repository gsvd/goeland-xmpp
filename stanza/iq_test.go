package stanza

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/gsvd/goeland-xmpp/internal/id"
	"github.com/stretchr/testify/require"
)

func TestIQSuite(t *testing.T) {
	t.Run("new", testIQNew)
	t.Run("types", testIQTypes)
	t.Run("marshal", testIQMarshal)
	t.Run("unmarshal", testIQUnmarshal)
}

func testIQNew(t *testing.T) {
	t.Run("valid type", func(t *testing.T) {
		iq, err := NewIQSet(
			WithIQLang("en"),
			WithIQFrom("alice@example.com/desktop"),
			WithIQTo("hector@example.com/tablet"),
			WithBindResource("desktop"),
			WithBindAddressStr("192.168.1.1"),
		)
		require.NoError(t, err)
		require.NotNil(t, iq)

		require.NotEmpty(t, iq.ID, "id should be generated")
		require.Equal(t, IQTypeSet, iq.Type)
		require.Equal(t, "en", iq.Lang)
		require.Equal(t, "alice@example.com/desktop", iq.From)
		require.Equal(t, "hector@example.com/tablet", iq.To)

		require.NotNil(t, iq.Bind)
		require.Equal(t, "desktop", iq.Bind.Resource)
		require.Equal(t, "192.168.1.1", iq.Bind.Address)

	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := NewIQ(IQType("bad-type"))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidIQType)
	})
}

func testIQTypes(t *testing.T) {
	tests := []struct {
		tc       string
		iqType   IQType
		expected string
	}{
		{"type get", IQTypeGet, "get"},
		{"type set", IQTypeSet, "set"},
		{"type result", IQTypeResult, "result"},
		{"type error", IQTypeError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			require.Equal(t, tt.expected, string(tt.iqType))
		})
	}
}

func testIQMarshal(t *testing.T) {
	runUUID := id.New()

	tests := []struct {
		tc       string
		input    IQ
		expected string
	}{
		{
			tc: "bind set with resource",
			input: IQ{
				Lang: "en",
				ID:   fmt.Sprintf("%s-iq0", runUUID),
				Type: IQTypeSet,
				From: "alice@example.com/desktop",
				To:   "example.com",
				Bind: &Bind{Resource: "desktop"},
			},
			expected: `<iq xml:lang="en" id="` + runUUID + `-iq0" type="set" from="alice@example.com/desktop" to="example.com"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><resource>desktop</resource></bind></iq>`,
		},
		{
			tc: "bind set without resource",
			input: IQ{
				ID:   fmt.Sprintf("%s-iq1", runUUID),
				Type: IQTypeSet,
				Bind: &Bind{},
			},
			expected: `<iq id="` + runUUID + `-iq1" type="set"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"></bind></iq>`,
		},
		{
			tc: "bind result with jid",
			input: IQ{
				ID:   fmt.Sprintf("%s-iq2", runUUID),
				Type: IQTypeResult,
				Bind: &Bind{Address: "alice@example.com/desktop"},
			},
			expected: `<iq id="` + runUUID + `-iq2" type="result"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><jid>alice@example.com/desktop</jid></bind></iq>`,
		},
		{
			tc: "error iq without child",
			input: IQ{
				ID:   fmt.Sprintf("%s-iq3", runUUID),
				Type: IQTypeError,
				From: "example.com",
			},
			expected: `<iq id="` + runUUID + `-iq3" type="error" from="example.com"></iq>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			data, err := xml.Marshal(&tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(data))
		})
	}
}

func testIQUnmarshal(t *testing.T) {
	runUUID := id.New()

	tests := []struct {
		tc       string
		input    string
		expected IQ
	}{
		{
			tc:    "bind set with resource",
			input: `<iq xml:lang="en" id="` + runUUID + `-iq0" type="set" from="alice@example.com/desktop" to="example.com"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><resource>desktop</resource></bind></iq>`,
			expected: IQ{
				XMLName: xml.Name{Local: "iq"},
				Lang:    "en",
				ID:      fmt.Sprintf("%s-iq0", runUUID),
				Type:    IQTypeSet,
				From:    "alice@example.com/desktop",
				To:      "example.com",
				Bind: &Bind{
					XMLName:  xml.Name{Space: NSBind, Local: "bind"},
					Resource: "desktop",
				},
			},
		},
		{
			tc:    "bind set without resource",
			input: `<iq id="` + runUUID + `-iq1" type="set"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"></bind></iq>`,
			expected: IQ{
				XMLName: xml.Name{Local: "iq"},
				ID:      fmt.Sprintf("%s-iq1", runUUID),
				Type:    IQTypeSet,
				Bind: &Bind{
					XMLName: xml.Name{Space: NSBind, Local: "bind"},
				},
			},
		},
		{
			tc:    "bind result with address",
			input: `<iq id="` + runUUID + `-iq2" type="result"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><jid>alice@example.com/desktop</jid></bind></iq>`,
			expected: IQ{
				XMLName: xml.Name{Local: "iq"},
				ID:      fmt.Sprintf("%s-iq2", runUUID),
				Type:    IQTypeResult,
				Bind: &Bind{
					XMLName: xml.Name{Space: NSBind, Local: "bind"},
					Address: "alice@example.com/desktop",
				},
			},
		},
		{
			tc:    "error iq without child",
			input: `<iq id="` + runUUID + `-iq3" type="error" from="example.com"></iq>`,
			expected: IQ{
				XMLName: xml.Name{Local: "iq"},
				ID:      fmt.Sprintf("%s-iq3", runUUID),
				Type:    IQTypeError,
				From:    "example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			var iq IQ
			err := xml.Unmarshal([]byte(tt.input), &iq)
			require.NoError(t, err)
			require.Equal(t, tt.expected, iq)
		})
	}
}
