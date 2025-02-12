package instance

import (
	"errors"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image/color"
	"sync"
	"sync/atomic"
)

type Player struct {
	*player.Player

	instance   Instance
	instanceMu sync.Mutex

	c session.Conn

	transferring atomic.Bool
}

func (pl *Player) Transferring() bool {
	return pl.transferring.Load()
}

func (pl *Player) setTransferring(t bool) {
	pl.transferring.Store(t)
}

func (pl *Player) setInstance(i Instance) {
	if i == nil {
		i = Nop
	}

	pl.instanceMu.Lock()
	defer pl.instanceMu.Unlock()
	pl.instance = i
}

func (pl *Player) Instance() Instance {
	pl.instanceMu.Lock()
	defer pl.instanceMu.Unlock()

	if pl.instance == nil {
		return Nop
	}

	return pl.instance
}

func GetTypedInstance[T any](pl *Player) (T, bool) {
	for _, i := range AllInstances() {
		if i.Active(pl.UUID()) {
			v, ok := i.(T)
			return v, ok
		}
	}

	var nop T
	return nop, false
}

type WithinTransaction func(p *player.Player)

func (pl *Player) ExecSafe(f func(*player.Player, *world.Tx), w ...WithinTransaction) {
	go pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
		f(e.(*player.Player), tx)

		go func() {
			for _, w := range w {
				w(e.(*player.Player))
			}
		}()
	})
}

func (pl *Player) SendKit(k kit.Kit, tx *world.Tx) {
	pl.Reset(tx, func(*player.Player) {
		pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
			kit.Send(k, p)
		})
	})
}

func (pl *Player) Conn() (session.Conn, bool) {
	return pl.c, pl.c != nil
}

func (pl *Player) MustConn() session.Conn {
	c, ok := pl.Conn()
	if !ok {
		Kick(pl.Player, ErrorSponge)
		return nil
	}

	return c
}

func (pl *Player) limitChunks(v int) {
	if c := pl.MustConn(); c != nil {
		_ = c.WritePacket(&packet.ChunkRadiusUpdated{ChunkRadius: int32(v)})
	}
}

func (pl *Player) HideElements(e ...hud.Element) error {
	if conn, ok := pl.Conn(); ok {
		hud.Hide(conn, e...)
		return nil
	}

	return errors.New("cannot get player connection")
}

func resetFunctions(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()

	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}

	p.SetGameMode(world.GameModeSurvival)

	p.CloseForm()
	p.CloseDialogue()

	p.EnableInstantRespawn()
	p.SetMobile()

	p.SetVisible()
	p.SetScale(1.0)

	p.ShowCoordinates()
}

func (pl *Player) Reset(tx *world.Tx, after ...func(p *player.Player)) {
	selfReset := func(pl *Player) {
		pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
			resetFunctions(p)

			for _, f := range after {
				f(p)
			}
		})
	}

	if tx == nil {
		selfReset(pl)
		return
	}

	e, ok := pl.H().Entity(tx)
	if !ok {
		selfReset(pl)
		return
	}

	p := e.(*player.Player)
	resetFunctions(p)

	for _, f := range after {
		f(p)
	}
}

func (pl *Player) World() *world.World {
	var w *world.World

	pl.ExecSafe(func(player *player.Player, tx *world.Tx) {
		w = tx.World()
	})

	return w
}

func (pl *Player) ResetElements(e ...hud.Element) error {
	if conn, ok := pl.Conn(); ok {
		hud.Reset(conn, e...)
		return nil
	}

	return errors.New("cannot get player connection")
}

func (pl *Player) syncConn(m *plugin.Manager) error {
	c, ok := m.Conn(pl.Name())
	if !ok {
		return errors.New("could not find connection")
	}

	pl.c = c
	return nil
}

func (pl *Player) enableChunkCache() {
	if c, ok := pl.Conn(); ok {
		_ = c.WritePacket(&packet.ClientCacheStatus{
			Enabled: true,
		})
	}
}

func FadeInCamera(c session.Conn, dur float32, fadeIn bool) {
	var duration, fadeInDuration float32

	duration = dur / 2

	if fadeIn {
		fadeInDuration = dur / 3
		duration = dur / 3
	}

	_ = c.WritePacket(&packet.CameraInstruction{
		Fade: protocol.Option(protocol.CameraInstructionFade{
			TimeData: protocol.Option(protocol.CameraFadeTimeData{
				FadeInDuration:  fadeInDuration,
				WaitDuration:    duration,
				FadeOutDuration: duration,
			}),
			Colour: protocol.Option(color.RGBA{}),
		}),
	})
}
