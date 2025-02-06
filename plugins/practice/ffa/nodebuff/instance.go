package nodebuff

import (
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance"
)

const name = "nodebuff"

func Instance() *ffa.Instance {
	return instance.MustByName(name).(*ffa.Instance)
}
