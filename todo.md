# Progress

## v0.1.0

### Features done

- [x] JID parsing and validation (RFC 7622)
- [x] Message stanza (basic send/receive, marshal/unmarshal) (RFC 6121, RFC 6120)
- [x] Message type support: chat, groupchat, headline, normal, error
- [x] Presence stanza (send/receive, marshal/unmarshal, types, show, status, priority)
- [x] IQ stanza (bind, ping, version, types, marshal/unmarshal)
- [x] Unit tests for address, message, presence, and IQ stanzas

### Features to implement

- [ ] XMPP stream handling (open/close XML streams, <stream:stream> root)
- [ ] TCP connection to XMPP server (connect, send/receive XML)
- [ ] Basic authentication (SASL PLAIN or ANONYMOUS)
- [ ] Basic event loop (read/write stanzas, dispatch handlers)
- [ ] Error handling for stream and stanza parsing
- [ ] Minimal example: connect, authenticate, send/receive message, presence

---

## Estimated Progress

- **Current completion:** 40% for a minimal XMPP chat client (v0.1)
- **Notes:** JID, message, presence, and IQ stanzas are implemented with tests. Stream handling, connection, authentication, and advanced IQ features are next.

---

_This file is auto-generated from .copilot/.copilot-todo.yaml. Update the YAML file as features are completed._
