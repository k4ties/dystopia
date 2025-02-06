package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
)

type FFA struct {
	onlyPlayer
}

func (c FFA) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	if !lobby.Instance().Active(inPl(s).UUID()) {
		o.Errorf("Can only teleport to FFA in lobby")
		return
	}

	p(s).SendForm(ffa.NewForm())
}

func init() {
	cmd.Register(cmd.New("ffa", "", nil, FFA{}))
}
