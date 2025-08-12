# Progress

## v0.1.0

### Features done

- [x] JID parsing and validation (RFC 7622)
- [x] Message stanza (basic send/receive, marshal/unmarshal) (RFC 6121, RFC 6120)
- [x] Message type support: chat, groupchat, headline, normal, error
- [x] Unit tests for address parsing and message stanzas

### Features to implement

- [ ] XMPP stream handling (open/close XML streams, <stream:stream> root)
- [ ] TCP connection to XMPP server (connect, send/receive XML)
- [ ] Basic authentication (SASL PLAIN or ANONYMOUS)
- [ ] Presence stanza (send/receive <presence/>)
- [ ] IQ stanza (basic, e.g., ping, version)
- [ ] Basic event loop (read/write stanzas, dispatch handlers)
- [ ] Error handling for stream and stanza parsing
- [ ] Minimal example: connect, authenticate, send/receive message, presence

---

## Estimated Progress

- **Current completion:** 20% for a minimal XMPP chat client (v0.1)
- **Notes:** JID and message stanzas are implemented. Stream, connection, authentication, presence, and IQ are next.

---

_This file is auto-generated from .copilot/.copilot-todo.yaml. Update the YAML file as features are completed._
