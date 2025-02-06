package lobby

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

const name = "lobby"

type impl struct {
	*instance.Impl
}

func Instance() instance.Instance {
	return &impl{Impl: instance.MustByName(name).(*instance.Impl)}
}

func (i *impl) Transfer(p *instance.Player, tx *world.Tx) {
	i.Impl.Transfer(p, tx)
}
