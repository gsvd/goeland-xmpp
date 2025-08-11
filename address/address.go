package address

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

var (
	ErrInvalidAddressFormat = errors.New("invalid address format")
	ErrAddressIsEmpty       = errors.New("address is empty")
	ErrPartInvalidUTF8      = errors.New("part contains invalid UTF-8")
	ErrPartLenIsTooLong     = errors.New("too long max 1023 bytes")
	ErrNotAllowedCharacters = errors.New("not allowed characters")

	ErrMissingResourcePart = errors.New("missing resource part")
	ErrMissingLocalPart    = errors.New("missing local part")
	ErrMissingDomainPart   = errors.New("missing domain part")
)

type Address struct {
	// local identifier (optional)
	local string
	// domain identifier (required)
	domain string
	// resource identifier (optional)
	resource string
}

func new(local, domain, resource string) (Address, error) {
	local, err := normalizeLocal(local)
	if err != nil {
		return Address{}, fmt.Errorf("invalid local part: %w", err)
	}

	domain, err = normalizeDomain(domain)
	if err != nil {
		return Address{}, fmt.Errorf("invalid domain part: %w", err)
	}

	resource, err = normalizeResource(resource)
	if err != nil {
		return Address{}, fmt.Errorf("invalid resource part: %w", err)
	}

	return Address{
		local:    local,
		domain:   domain,
		resource: resource,
	}, nil
}

func Parse(addr string) (Address, error) {
	local, domain, resource, err := decompose(addr)
	if err != nil {
		return Address{}, err
	}

	address, err := new(local, domain, resource)
	if err != nil {
		return Address{}, err
	}

	return address, nil
}

func MustParse(addr string) Address {
	local, domain, resource, err := decompose(addr)
	if err != nil {
		panic(err)
	}

	address, err := new(local, domain, resource)
	if err != nil {
		panic(err)
	}

	return address
}

func decompose(addr string) (local string, domain string, resource string, err error) {
	if addr == "" {
		err = ErrAddressIsEmpty
		return
	}

	sep := strings.LastIndex(addr, "/")
	if sep > -1 {
		// Trailing slash is present, but resource is missing
		// e.g. user@example.com/
		if sep == len(addr)-1 {
			err = ErrMissingResourcePart
			return
		}
		resource = addr[sep+1:]
		addr = addr[:sep]
	}

	sep = strings.LastIndex(addr, "@")
	switch {
	case sep == -1:
		// No @, entire string is domain
		domain = addr
	case sep == len(addr)-1:
		// @ at end, e.g. user@
		err = ErrMissingDomainPart
		return
	case sep == 0:
		// @ at start, e.g. @example.com
		err = ErrMissingLocalPart
		return
	default:
		domain = addr[sep+1:]
		local = addr[:sep]
	}

	return
}

// TODO: Implement
func normalizeLocal(local string) (string, error) {
	if len(local) > 1023 {
		return local, ErrPartLenIsTooLong
	}

	if !utf8.ValidString(local) {
		return local, ErrPartInvalidUTF8
	}

	if strings.ContainsFunc(local, func(r rune) bool {
		switch r {
		case '\u0022', // " QUOTATION MARK
			'\u0026', // & AMPERSAND
			'\u0027', // ' APOSTROPHE
			'\u002F', // / SOLIDUS
			'\u003A', // : COLON
			'\u003C', // < LESS-THAN SIGN
			'\u003E', // > GREATER-THAN SIGN
			'\u0040': // @ COMMERCIAL AT
			return true
		default:
			return false
		}
	}) {
		return local, ErrNotAllowedCharacters
	}

	return local, nil
}

func normalizeDomain(domain string) (string, error) {
	if !utf8.ValidString(domain) {
		return domain, ErrPartInvalidUTF8
	}

	if strings.HasPrefix(domain, "[") && strings.HasSuffix(domain, "]") {
		ipv6 := domain[1 : len(domain)-1]
		ip := net.ParseIP(ipv6)
		if ip != nil && ip.To4() == nil {
			return domain, nil
		}
	}

	if ip := net.ParseIP(domain); ip != nil && ip.To4() != nil {
		return domain, nil
	}

	domain = strings.TrimSuffix(domain, ".")

	domain, err := idna.Display.ToUnicode(domain)
	if err != nil {
		return domain, err
	}

	if len(domain) > 1023 {
		return domain, ErrPartLenIsTooLong
	}

	return domain, nil
}

// TODO: Implement
func normalizeResource(resource string) (string, error) {
	if len(resource) > 1023 {
		return resource, ErrPartLenIsTooLong
	}

	return resource, nil
}

func (a Address) Equal(b Address) bool {
	return a.local == b.local && a.domain == b.domain && a.resource == b.resource
}

func (a Address) String() string {
	var b strings.Builder
	b.Grow(len(a.local) + len(a.domain) + len(a.resource) + 2)

	if a.local != "" {
		b.WriteString(a.local)
	}

	if a.local != "" && a.domain != "" {
		b.WriteByte('@')
	}

	b.WriteString(a.domain)

	if a.resource != "" {
		b.WriteByte('/')
		b.WriteString(a.resource)
	}

	return b.String()
}

func (a Address) Bare() Address {
	return Address{
		local:  a.local,
		domain: a.domain,
	}
}

func (a Address) Local() Address {
	return Address{
		local: a.local,
	}
}

func (a Address) Domain() Address {
	return Address{
		domain: a.domain,
	}
}
