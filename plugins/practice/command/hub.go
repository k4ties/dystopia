package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
)

type Hub struct {
	onlyPlayer
}

func (Hub) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	lobby.TransferWithRoutine(inPl(s), tx)
	systemMessage(o, "You've been teleported to the Lobby.")
}
