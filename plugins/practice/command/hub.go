package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
)

type Hub struct {
	onlyPlayer
}

func (c Hub) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	lobby.AddToWorld(instance.LookupPlayer(s.(*player.Player)))

	systemMessage(s, "You've been teleported to the hub.")
}
