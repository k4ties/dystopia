package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

type ReKit struct {
	onlyPlayer
}

func (ReKit) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	if pl := inPl(s); pl != nil {
		for _, i := range instance.AllInstances() {
			if f, ok := i.(*ffa.Instance); ok {
				if f.Active(pl.UUID()) {
					pl.SendKit(f.Kit(), tx)
					systemMessage(o, "You've successfully rekitted.")
					return
				}
			}
		}

		o.Errorf("You're not in ffa.")
	}
}
