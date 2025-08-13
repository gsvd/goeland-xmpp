package address

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

type Address struct {
	local    string
	domain   string
	resource string
}

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

type Option func(*Address)

func New(opts ...Option) (*Address, error) {
	a := &Address{}

	for _, opt := range opts {
		opt(a)
	}

	if a.domain == "" {
		return nil, ErrMissingDomainPart
	}

	if a.local != "" {
		l, err := normalizeLocal(a.local)
		if err != nil {
			return nil, fmt.Errorf("invalid local part: %w", err)
		}
		a.local = l
	}

	d, err := normalizeDomain(a.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid domain part: %w", err)
	}
	a.domain = d

	if a.resource != "" {
		r, err := normalizeResource(a.resource)
		if err != nil {
			return nil, fmt.Errorf("invalid resource part: %w", err)
		}
		a.resource = r
	}

	return a, nil
}

func Parse(s string) (*Address, error) {
	var opts []Option

	local, domain, resource, err := decompose(s)
	if err != nil {
		return nil, err
	}

	if local != "" {
		opts = append(opts, WithLocal(local))
	}

	opts = append(opts, WithDomain(domain))

	if resource != "" {
		opts = append(opts, WithResource(resource))
	}

	return New(opts...)
}

func MustParse(addr string) *Address {
	address, err := Parse(addr)
	if err != nil {
		panic(err)
	}

	return address
}

func WithLocal(local string) Option {
	return func(a *Address) {
		a.local = local
	}
}

func WithDomain(domain string) Option {
	return func(a *Address) {
		a.domain = domain
	}
}

func WithResource(resource string) Option {
	return func(a *Address) {
		a.resource = resource
	}
}

func decompose(addr string) (local string, domain string, resource string, err error) {
	if addr == "" {
		err = ErrAddressIsEmpty
		return
	}

	sep := strings.LastIndex(addr, "/")
	if sep > -1 {
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
		domain = addr
	case sep == len(addr)-1:
		err = ErrMissingDomainPart
		return
	case sep == 0:
		err = ErrMissingLocalPart
		return
	default:
		domain = addr[sep+1:]
		local = addr[:sep]
	}

	return
}

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

	d, err := idna.Display.ToUnicode(domain)
	if err != nil {
		return domain, err
	}

	if len(d) > 1023 {
		return d, ErrPartLenIsTooLong
	}

	return d, nil
}

func normalizeResource(resource string) (string, error) {
	if len(resource) > 1023 {
		return resource, ErrPartLenIsTooLong
	}

	if !utf8.ValidString(resource) {
		return resource, ErrPartInvalidUTF8
	}

	return resource, nil
}

func (a *Address) Equal(b *Address) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.local == b.local && a.domain == b.domain && a.resource == b.resource
}

func (a *Address) String() string {
	if a == nil {
		return ""
	}

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

func (a *Address) Bare() *Address {
	if a == nil {
		return nil
	}
	return &Address{
		local:  a.local,
		domain: a.domain,
	}
}

func (a *Address) Local() *Address {
	if a == nil {
		return nil
	}
	return &Address{
		local: a.local,
	}
}

func (a *Address) Domain() *Address {
	if a == nil {
		return nil
	}
	return &Address{
		domain: a.domain,
	}
}
