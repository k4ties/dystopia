package ffa

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"slices"
	"strings"
	"sync/atomic"
	"time"
)

type Instance struct {
	*instance.Impl
	k kit.Kit

	c      Config
	closed atomic.Bool
}

type Config struct {
	Name string
	Icon string

	PearlCooldown time.Duration
}

func New(i *instance.Impl, k kit.Kit, c Config) *Instance {
	f := &Instance{k: k, Impl: i, c: c}

	if slices.Contains(Closed.Closed, strings.ToLower(c.Name)) {
		i.World().Exec(func(tx *world.Tx) {
			f.Close(tx)
		})
	}

	return f
}

func (i *Instance) HasPearCooldown() bool {
	return i.c.PearlCooldown > 0
}

func (i *Instance) PearlCooldown() time.Duration {
	return i.c.PearlCooldown
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

	pl.Messagef(text.Colourf("<green>You've been teleported to the %s.</green>", i.c.Name))
}

func (i *Instance) RemoveFromList(pl *instance.Player) {
	i.Impl.RemoveFromList(pl)
	i.Messagef("<red>%s</red> left from the <red>%s</red>. <dark-grey>(%d)</dark-grey>", pl.Name(), i.c.Name, i.playerLen())
}

func (i *Instance) playerLen() int {
	var l int

	for range i.Players() {
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
		p.Messagef(text.Colourf("<red>This game mode has been closed.</red>"))
		lobby.TransferWithRoutine(p, tx)
	}
}

func (i *Instance) Name() string {
	return i.c.Name
}

func (i *Instance) Icon() string {
	return i.c.Icon
}
