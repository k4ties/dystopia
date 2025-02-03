package lobby

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

func i() instance.Instance {
	return Instance()
}

func Instance() instance.Instance {
	return instance.MustByName("lobby")
}

func Player(p *player.Player) *instance.Player {
	if !i().Active(p.UUID()) {
		return nil
	}

	return i().NewPlayer(p)
}

func AddToWorld(p *instance.Player) {
	if p == nil {
		panic("instance player cannot be nil")
	}

	if !i().Active(p.UUID()) {
		if p.Instance() != nil {
			p.Instance().RemovePlayer(p)
		}

		i().AddPlayer(p)
		return
	}

	p.ExecSafe(func(pl *player.Player, tx *world.Tx) {
		pl.Teleport(i().World().Spawn().Vec3Centre())
	})
}
