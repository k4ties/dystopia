package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type Gm struct {
	onlyPlayer
	// no arguments
}

func (Gm) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := p(s)

	switch p.GameMode() {
	case world.GameModeCreative:
		p.SetGameMode(world.GameModeSurvival)
	default:
		p.SetGameMode(world.GameModeCreative)
	}
}

func init() {
	cmd.Register(cmd.New("gm", "", nil, Gm{}))
}
