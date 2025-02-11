package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
)

type ReKit struct {
	onlyPlayer
}

func (ReKit) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, in := ffa.LookupPlayer(s.(*player.Player))
	if pl == nil || in == nil {
		o.Errorf("You're not in ffa.")
		return
	}

	in.ReKit(pl, tx)
	pl.Heal(pl.MaxHealth(), effect.InstantHealingSource{})

	systemMessage(o, "You've successfully rekitted.")
}
