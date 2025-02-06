package lobby

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

const name = "lobby"

func Rotation() cube.Rotation {
	return cube.Rotation{180, 0}
}

func Instance() instance.Instance {
	return instance.MustByName(name)
}
