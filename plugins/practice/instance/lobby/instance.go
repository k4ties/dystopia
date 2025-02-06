package lobby

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

const name = "lobby"

func Instance() instance.Instance {
	return instance.MustByName(name)
}

func TransferWithRoutine(pl *instance.Player, tx *world.Tx) {
	Instance().Transfer(pl, tx)

	pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
		if Instance().Active(pl.UUID()) {
			p.Teleport(Instance().World().Spawn().Vec3Centre())
		}

		pl.SendKit(Kit, tx)
	})
}
