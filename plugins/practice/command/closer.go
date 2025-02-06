package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"strings"
)

type Closer struct {
	onlyPlayer
	Instance string
}

func (c Closer) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	i, ok := instance.ByName(strings.ToLower(c.Instance))
	if !ok {
		o.Errorf("No instances with name %s", c.Instance)
		return
	}

	f, isFfa := i.(*ffa.Instance)
	if !isFfa {
		o.Errorf("Instance %s is not a ffa.Instance, we cannot close it.", c.Instance)
		return
	}

	if f.Closed() {
		f.Open()
		o.Printf("Opened %s", c.Instance)
		return
	}

	f.Close(tx)
	o.Printf("Closed %s", c.Instance)
}

func init() {
	cmd.Register(cmd.New("close", "", nil, Closer{}))
}
