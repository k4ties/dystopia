package nodebuff

import (
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

const name = "nodebuff"

func Instance() instance.Instance {
	return instance.MustByName(name)
}
