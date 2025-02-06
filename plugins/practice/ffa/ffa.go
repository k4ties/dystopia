package ffa

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"sync/atomic"
)

type Instance struct {
	*instance.Impl
	k kit.Kit

	c Config

	closed atomic.Bool
}

type Config struct {
	Name string
	Icon string
}

func New(i *instance.Impl, k kit.Kit, c Config) *Instance {
	return &Instance{k: k, Impl: i, c: c}
}

func (i *Instance) Transfer(pl *instance.Player, tx *world.Tx) {
	if i.Closed() {
		pl.Messagef(text.Colourf("<red>Sorry, this game mode is closed.</red>"))
		return
	}

	if !i.Impl.Active(pl.UUID()) {
		i.Messagef("<red>%s</red> has joined to <red>%s</red>. <dark-grey>(%d)</dark-grey>", pl.Name(), i.c.Name, i.playerLen())
	}

	i.Impl.Transfer(pl, tx)
	pl.SendKit(i.k, tx)
}

func (i *Instance) RemoveFromList(pl *instance.Player) {
	i.Impl.RemoveFromList(pl)
	i.Messagef("<red>%s</red> left from the <red>%s</red>. <dark-grey>(%d)</dark-grey>", pl.Name(), i.c.Name, i.playerLen())
}

func (i *Instance) playerLen() int {
	var l int

	for _ = range i.Players() {
		l++
	}

	return l
}

func (i *Instance) Closed() bool {
	return i.closed.Load()
}

func (i *Instance) Open() {
	i.closed.Store(false)
}

func (i *Instance) Close(tx *world.Tx) {
	i.closed.Store(true)

	for p := range i.Impl.Players() {
		p.Messagef("<red>This game mode has been closed.</red>")
		lobby.Instance().Transfer(p, p.Tx())
	}
}

func (i *Instance) Name() string {
	return i.c.Name
}

func (i *Instance) Icon() string {
	return i.c.Icon
}
