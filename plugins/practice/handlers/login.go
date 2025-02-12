package handlers

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/session"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/handlers/whitelist"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strings"
)

type Login struct {
	plugin.NopPlayerHandler
}

func NewLoginHandler() *Login {
	l := &Login{}

	return l
}

func (l *Login) HandleLogin(ctx *event.Context[session.Conn]) {
	if whitelist.Enabled() {
		name := strings.ToLower(ctx.Val().IdentityData().DisplayName)

		if !whitelist.Whitelisted(name) {
			_ = ctx.Val().WritePacket(&packet.Disconnect{
				Message: "Server is whitelisted.",
			})
			ctx.Cancel()
			return
		}
	}
}
