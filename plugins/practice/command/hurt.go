package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Hurt struct {
	onlyPlayer
	Amount float64
}

func (c Hurt) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	a := c.Amount

	inPl(s).ExecSafe(func(p *player.Player, _ *world.Tx) {
		p.Hurt(a, entity.AttackDamageSource{})
	})
}

func init() {
	cmd.Register(cmd.New("hurt", "", nil, Hurt{}))
}
