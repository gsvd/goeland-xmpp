package address

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressSuite(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		testAddressParsing(t)
	})

	t.Run("stringify", func(t *testing.T) {
		testAddressStringify(t)
	})

	t.Run("local part invalid characters", func(t *testing.T) {
		testLocalPartInvalidCharacters(t)
	})
}

func testAddressParsing(t *testing.T) {
	tests := []struct {
		tc        string
		input     string
		expectRes Address
		expectErr error
	}{
		{
			tc:        "valid address",
			input:     "user@example.com/resource",
			expectRes: Address{local: "user", domain: "example.com", resource: "resource"},
			expectErr: nil,
		},
		{
			tc:        "missing resource part",
			input:     "user@example.com/",
			expectRes: Address{},
			expectErr: ErrMissingResourcePart,
		},
		{
			tc:        "missing local part",
			input:     "@example.com/resource",
			expectRes: Address{},
			expectErr: ErrMissingLocalPart,
		},
		{
			tc:        "empty address",
			input:     "",
			expectRes: Address{},
			expectErr: ErrAddressIsEmpty,
		},
		{
			tc:        "domain is too short",
			input:     "user@/resource",
			expectRes: Address{},
			expectErr: ErrMissingDomainPart,
		},
		{
			tc:        "domain is too long",
			input:     fmt.Sprintf("user@%s.com", strings.Repeat("a", 1024)),
			expectRes: Address{},
			expectErr: ErrPartLenIsTooLong,
		},
		{
			tc:        "just an @",
			input:     "@",
			expectRes: Address{},
			expectErr: ErrMissingDomainPart,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			res, err := Parse(tt.input)
			if tt.expectErr != nil {
				require.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectRes, res)
			}
		})
	}
}

func testAddressStringify(t *testing.T) {
	tests := []struct {
		tc    string
		input string
	}{
		{
			tc:    "full address",
			input: "user@example.com/resource",
		},
		{
			tc:    "bare address",
			input: "user@example.com",
		},
		{
			tc:    "domain only",
			input: "example.com",
		},
		{
			tc:    "domain with resource and no local",
			input: "example.com/resource",
		},
		{
			tc:    "ipv4 address as domain",
			input: "user@192.168.1.1/resource",
		},
		{
			tc:    "ipv4 address bare",
			input: "192.168.1.1",
		},
		{
			tc:    "ipv6 address as domain",
			input: "user@[2001:db8::1]/resource",
		},
		{
			tc:    "ipv6 address with resource only",
			input: "[2001:db8::1]/resource",
		},
		{
			tc:    "muc room participant",
			input: "conference.example.com/nickname",
		},
		{
			tc:    "subdomain only",
			input: "muc.conference.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			address := MustParse(tt.input)
			require.Equal(t, tt.input, address.String())
		})
	}
}

func testLocalPartInvalidCharacters(t *testing.T) {
	tests := []struct {
		tc        string
		input     string
		expectErr error
	}{
		{
			tc:        "invalid character quotation mark in local part",
			input:     `user"name@example.com/resource`,
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character ampersand in local part",
			input:     "user&name@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character apostrophe in local part",
			input:     "us'er@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character solidus in local part",
			input:     "user/name@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character colon in local part",
			input:     "user:name@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character less-than in local part",
			input:     "user<@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character greater-than in local part",
			input:     "user>name@example.com/resource",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character quotation mark at start of local part",
			input:     `"username@example.com`,
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "invalid character ampersand at end of local part",
			input:     "username&@example.com",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "multiple invalid characters in local part",
			input:     "user<name>@example.com",
			expectErr: ErrNotAllowedCharacters,
		},
		{
			tc:        "valid characters in local part",
			input:     "user.name@example.com/resource",
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			_, err := Parse(tt.input)
			if tt.expectErr != nil {
				require.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
