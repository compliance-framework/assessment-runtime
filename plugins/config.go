package plugins

import "github.com/hashicorp/go-plugin"

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "AR_PLUGIN",
	MagicCookieValue: "048cc450-6be2-4fa2-b760-8d4d0b63b534",
}
