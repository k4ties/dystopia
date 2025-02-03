package handlers

import (
	"github.com/df-mc/dragonfly/server/player"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

type Instance struct {
	plugin.NopPlayerHandler
	i instance.Instance
}

func NewInstance(i instance.Instance) *Instance {
	return &Instance{i: i}
}

func (l *Instance) HandleSpawn(p *player.Player) {
	l.i.AddPlayer(instance.NewPlayer(p))
}

func (l *Instance) HandleQuit(p *player.Player) {
	if pl := instance.LookupPlayer(p); pl != nil {
		pl.Instance().RemovePlayer(pl)
	}
}
