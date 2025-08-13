package address

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressSuite(t *testing.T) {
	t.Run("parsing", testAddressParsing)
	t.Run("stringify", testAddressStringify)
	t.Run("local part invalid characters", testLocalPartInvalidCharacters)
	t.Run("new with options", testAddressNewOptions)
	t.Run("helpers bare/local/domain", testAddressHelpers)
	t.Run("mustparse panics", testAddressMustParsePanics)
	t.Run("validation edges", testAddressValidationEdges)
}

func testAddressParsing(t *testing.T) {
	tests := []struct {
		tc        string
		input     string
		expectRes *Address
		expectErr error
	}{
		{
			tc:        "valid address",
			input:     "user@example.com/resource",
			expectRes: &Address{local: "user", domain: "example.com", resource: "resource"},
			expectErr: nil,
		},
		{
			tc:        "missing resource part",
			input:     "user@example.com/",
			expectRes: &Address{},
			expectErr: ErrMissingResourcePart,
		},
		{
			tc:        "missing local part",
			input:     "@example.com/resource",
			expectRes: &Address{},
			expectErr: ErrMissingLocalPart,
		},
		{
			tc:        "empty address",
			input:     "",
			expectRes: &Address{},
			expectErr: ErrAddressIsEmpty,
		},
		{
			tc:        "domain is too short",
			input:     "user@/resource",
			expectRes: &Address{},
			expectErr: ErrMissingDomainPart,
		},
		{
			tc:        "domain is too long",
			input:     fmt.Sprintf("user@%s.com", strings.Repeat("a", 1024)),
			expectRes: &Address{},
			expectErr: ErrPartLenIsTooLong,
		},
		{
			tc:        "just an @",
			input:     "@",
			expectRes: &Address{},
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

func testAddressNewOptions(t *testing.T) {
	tests := []struct {
		tc        string
		opts      []Option
		expectRes *Address
		expectErr error
	}{
		{
			tc: "domain only",
			opts: []Option{
				WithDomain("example.com"),
			},
			expectRes: &Address{domain: "example.com"},
			expectErr: nil,
		},
		{
			tc: "local and domain and resource",
			opts: []Option{
				WithLocal("user"),
				WithDomain("example.com"),
				WithResource("work"),
			},
			expectRes: &Address{local: "user", domain: "example.com", resource: "work"},
			expectErr: nil,
		},
		{
			tc:        "missing domain",
			opts:      []Option{WithLocal("user")},
			expectRes: &Address{},
			expectErr: ErrMissingDomainPart,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			res, err := New(tt.opts...)
			if tt.expectErr != nil {
				require.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectRes, res)
			}
		})
	}
}

func testAddressHelpers(t *testing.T) {
	addr := MustParse("user@example.com/work")

	t.Run("bare", func(t *testing.T) {
		require.Equal(t, &Address{local: "user", domain: "example.com"}, addr.Bare())
	})
	t.Run("local", func(t *testing.T) {
		require.Equal(t, &Address{local: "user"}, addr.Local())
	})
	t.Run("domain", func(t *testing.T) {
		require.Equal(t, &Address{domain: "example.com"}, addr.Domain())
	})
	t.Run("equal", func(t *testing.T) {
		a1 := MustParse("user@example.com/work")
		a2 := MustParse("user@example.com/work")
		require.True(t, a1.Equal(a2))
	})
}

func testAddressMustParsePanics(t *testing.T) {
	tests := []struct {
		tc    string
		input string
	}{
		{"empty", ""},
		{"trailing slash no resource", "user@example.com/"},
		{"missing domain", "user@/res"},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			require.Panics(t, func() { _ = MustParse(tt.input) })
		})
	}
}

func testAddressValidationEdges(t *testing.T) {
	t.Run("resource too long", func(t *testing.T) {
		longRes := strings.Repeat("r", 1024)
		_, err := Parse("example.com/" + longRes)
		require.ErrorIs(t, err, ErrPartLenIsTooLong)
	})

	t.Run("invalid utf8 in local", func(t *testing.T) {
		_, err := Parse("us\xffer@example.com/res")
		require.ErrorIs(t, err, ErrPartInvalidUTF8)
	})

	t.Run("domain trailing dot normalized", func(t *testing.T) {
		addr, err := Parse("user@example.com./res")
		require.NoError(t, err)
		require.Equal(t, "user@example.com/res", addr.String())
	})

	t.Run("ipv6 domain bracketed", func(t *testing.T) {
		addr, err := Parse("[2001:db8::1]/r")
		require.NoError(t, err)
		require.Equal(t, &Address{domain: "[2001:db8::1]", resource: "r"}, addr)
		require.Equal(t, "[2001:db8::1]/r", addr.String())
	})
}
