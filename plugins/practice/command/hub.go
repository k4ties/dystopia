package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Hub struct {
	onlyPlayer
	// no arguments
}

func (Hub) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := lobby.Instance()
	l.Transfer(inPl(s), tx)

	if tx.World() == l.World() {
		p(s).Teleport(l.World().Spawn().Vec3Centre())
	}

	o.Printf(text.Colourf("<green>You've been teleported to the Lobby.</green>"))
}

func init() {
	cmd.Register(cmd.New("hub", "", nil, Hub{}))
}
