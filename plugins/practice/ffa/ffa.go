package ffa

import (
	"context"
	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Instance struct {
	*instance.Impl
	k kit.Kit

	c      Config
	closed atomic.Bool

	cooldowns   map[uuid.UUID]*cooldown.CoolDown
	cooldownsMu sync.RWMutex

	cooldownCancelFuncs   map[uuid.UUID]context.CancelFunc
	cooldownCancelFuncsMu sync.RWMutex
}

type Config struct {
	Name string
	Icon string

	PearlCooldown time.Duration
}

func New(i *instance.Impl, k kit.Kit, c Config) *Instance {
	f := &Instance{k: k, Impl: i, c: c, cooldowns: make(map[uuid.UUID]*cooldown.CoolDown), cooldownCancelFuncs: make(map[uuid.UUID]context.CancelFunc)}
	f.Impl = f.Impl.WithOnExitFuncs(f.OnExit)

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

const (
	JoinFormat = "<red>%s</red> has joined to <red>%s</red>. <dark-grey>(%d)</dark-grey>"
	QuitFormat = "<red>%s</red> left from the <red>%s</red>. <dark-grey>(%v)</dark-grey>"
)

func (i *Instance) Transfer(pl *instance.Player, tx *world.Tx) {
	if i.Closed() {
		pl.Messagef(text.Colourf("<red>Sorry, this game mode is closed.</red>"))
		return
	}

	if !i.Impl.Active(pl.UUID()) {
		msg := text.Colourf(JoinFormat, pl.Name(), i.c.Name, i.playerLen())

		i.Messagef(msg)
		pl.Messagef(msg)

		i.addToCooldownList(pl.UUID())
	}

	i.Impl.Transfer(pl, tx)

	pl.SendKit(i.k, tx)
	pl.Messagef(text.Colourf("<green>You've been teleported to the %s.</green>", i.c.Name))
}

func (i *Instance) ReKit(pl *instance.Player, tx *world.Tx) {
	if !i.Active(pl.UUID()) {
		return
	}

	i.ResetPearlCooldown(pl)
	pl.Heal(pl.MaxHealth(), effect.InstantHealingSource{})

	if tx == nil {
		pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
			pl.SendKit(i.Kit(), tx)
		})
		return
	}

	pl.SendKit(i.Kit(), tx)
}

func (i *Instance) StartPearlCoolDown(p *instance.Player, c context.Context) context.CancelFunc {
	if !i.Active(p.UUID()) {
		return nil
	}
	if c == nil {
		c = context.Background()
	}

	ctx, cancel := context.WithCancel(c)

	i.MustCoolDown(p.UUID(), func(cd *cooldown.CoolDown) {
		if cd.Active() {
			// if cooldown is already active we are resetting the existing one
			i.ResetPearlCooldown(p)
		}

		cd.Set(i.PearlCooldown())
		i.addCooldownCancelFunc(p.UUID(), cancel)

		p.SetExperienceLevel(int(i.PearlCooldown().Seconds()) - 1)
		p.SetExperienceProgress(0.99) // if we will make 1.0 it will be new level

		go func() {
			var showMessage = true

			defer func() {
				i.ResetPearlCooldown(p)

				if showMessage {
					p.Messagef(text.Colourf("<green>Your cooldown has been expired</green>"))
				}
			}()

			ticker := time.NewTicker(time.Second)
			timer := time.NewTimer(i.PearlCooldown() - time.Second)

			defer func() {
				ticker.Stop()
				timer.Stop()
				cd.Set(0)
			}()

			for {
				select {
				case <-ctx.Done():
					// if canceled, we must not show message
					showMessage = false
					return
				case <-timer.C:
					return
				case <-ticker.C:
					a := 1.0                                       // max value
					b := float64(int(i.PearlCooldown().Seconds())) // cooldown
					c := int(cd.Remaining().Seconds())             // step

					p.SetExperienceLevel(c)
					p.SetExperienceProgress((a / b) * float64(c))
				}
			}
		}()
	})

	return cancel
}

func (i *Instance) addCooldownCancelFunc(u uuid.UUID, c context.CancelFunc) {
	i.cooldownCancelFuncsMu.Lock()
	defer i.cooldownCancelFuncsMu.Unlock()

	i.cooldownCancelFuncs[u] = c
}

func (i *Instance) removeCooldownCancelFunc(u uuid.UUID) {
	i.cooldownCancelFuncsMu.Lock()
	defer i.cooldownCancelFuncsMu.Unlock()

	delete(i.cooldownCancelFuncs, u)
}

func (i *Instance) cooldownCancelFuncAvailable(u uuid.UUID) bool {
	i.cooldownCancelFuncsMu.RLock()
	defer i.cooldownCancelFuncsMu.RUnlock()

	return i.cooldownCancelFuncs[u] != nil
}

func (i *Instance) cooldownCancelFunc(u uuid.UUID) context.CancelFunc {
	i.cooldownCancelFuncsMu.RLock()
	defer i.cooldownCancelFuncsMu.RUnlock()

	return i.cooldownCancelFuncs[u]
}

func (i *Instance) OnExit(pl *instance.Player, _ instance.Instance) {
	i.ResetPearlCooldown(pl)
	i.removeFromCooldownList(pl.UUID())

	i.Messagef(text.Colourf(QuitFormat, pl.Name(), i.c.Name, i.playerLen()-1))
}

func (i *Instance) ResetPearlCooldown(pl *instance.Player) {
	if i.cooldownCancelFuncAvailable(pl.UUID()) {
		f := i.cooldownCancelFunc(pl.UUID())
		f()

		i.removeCooldownCancelFunc(pl.UUID())

		pl.SetExperienceLevel(0)
		pl.SetExperienceProgress(0)
	}
}

func (i *Instance) playerLen() int {
	var l int

	for range i.Players() {
		l++
	}

	return l + 1
}

func (i *Instance) CoolDown(u uuid.UUID) (*cooldown.CoolDown, bool) {
	i.cooldownsMu.RLock()
	defer i.cooldownsMu.RUnlock()

	cd, ok := i.cooldowns[u]
	return cd, ok
}

func (i *Instance) addToCooldownList(u uuid.UUID) {
	i.cooldownsMu.Lock()
	defer i.cooldownsMu.Unlock()

	i.cooldowns[u] = cooldown.NewCoolDown()
}

func (i *Instance) removeFromCooldownList(u uuid.UUID) {
	i.cooldownsMu.Lock()
	defer i.cooldownsMu.Unlock()

	delete(i.cooldowns, u)
}

func (i *Instance) MustCoolDown(u uuid.UUID, ifExists func(c *cooldown.CoolDown)) {
	c, ok := i.CoolDown(u)
	if !ok {
		return
	}

	ifExists(c)
}

func (i *Instance) Kit() kit.Kit {
	return i.k
}

func (i *Instance) Closed() bool {
	return i.closed.Load()
}

func (i *Instance) Open() {
	i.closed.Store(false)
}

func (i *Instance) Close(tx *world.Tx) {
	i.closed.Store(true)

	for p := range i.Players() {
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

func LookupPlayer(p *player.Player) (*instance.Player, *Instance) {
	for _, i := range instance.AllInstances() {
		if in, isFFA := i.(*Instance); isFFA {
			if in.Active(p.UUID()) {
				for pl := range in.Players() {
					if pl.UUID() == p.UUID() {
						return pl, in
					}
				}
				break
			}
		}
	}

	return nil, nil
}
