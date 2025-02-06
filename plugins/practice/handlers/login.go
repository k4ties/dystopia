package handlers

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/session"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"sync"
)

type Login struct {
	plugin.NopPlayerHandler

	whitelisted atomic.Bool
	players     sync.Map
}

func (l *Login) AddPlayer(pl string) {
	l.players.Store(pl, nil)
}

func (l *Login) DeletePlayer(pl string) {
	l.players.Delete(pl)
}

func (l *Login) ToggleWhitelist() {
	l.whitelisted.Store(!l.whitelisted.Load())
}

func NewLoginHandler(whitelisted bool, players []string) plugin.PlayerHandler {
	l := &Login{}
	l.whitelisted.Store(whitelisted)

	for _, pl := range players {
		l.players.Store(pl, nil)
	}

	return l
}

func (l *Login) HandleLogin(ctx *event.Context[session.Conn]) {
	if l.whitelisted.Load() {
		name := ctx.Val().IdentityData().DisplayName

		if _, ok := l.players.Load(name); !ok { // if we cannot load player from map, then this player is not whitelisted
			_ = ctx.Val().WritePacket(&packet.Disconnect{
				Message: "Server is whitelisted.",
			})
			ctx.Cancel()
			return
		}
	}
}
