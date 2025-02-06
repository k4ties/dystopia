package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/ffa/nodebuff"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Arena struct {
	onlyPlayer
	// no arguments
}

func (Arena) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	if dead(s) {
		o.Errorf("Cannot teleport while you're dead")
		return
	}

	n := nodebuff.Instance()
	n.Transfer(inPl(s), tx)

	if tx.World() == n.World() {
		p(s).Teleport(n.World().Spawn().Vec3Centre())
	}

	o.Printf(text.Colourf("<green>You've been teleported to the NoDebuff.</green>"))
}

func init() {
	cmd.Register(cmd.New("arena", "", nil, Arena{}))
}
