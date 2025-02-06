package instance

import (
	"errors"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image/color"
	"sync"
)

type Player struct {
	*player.Player

	instance   Instance
	instanceMu sync.Mutex

	c session.Conn
}

func (p *Player) setInstance(i Instance) {
	if i == nil {
		i = Nop
	}

	p.instanceMu.Lock()
	defer p.instanceMu.Unlock()
	p.instance = i
}

func (p *Player) Instance() Instance {
	p.instanceMu.Lock()
	defer p.instanceMu.Unlock()

	if p.instance == nil {
		return Nop
	}

	return p.instance
}

func (p *Player) ExecSafe(f func(*player.Player, *world.Tx)) {
	go p.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
		f(e.(*player.Player), tx)
	})
}

func (p *Player) Conn() (session.Conn, bool) {
	return p.c, p.c != nil
}

func (p *Player) HideElements(e ...hud.Element) error {
	if conn, ok := p.Conn(); ok {
		hud.Hide(conn, e...)
		return nil
	}

	return errors.New("cannot get player connection")
}

func (p *Player) World() *world.World {
	var w *world.World

	p.ExecSafe(func(player *player.Player, tx *world.Tx) {
		w = tx.World()
	})

	return w
}

func (p *Player) ResetElements(e ...hud.Element) error {
	if conn, ok := p.Conn(); ok {
		hud.Reset(conn, e...)
		return nil
	}

	return errors.New("cannot get player connection")
}

func (p *Player) syncConn(m *plugin.Manager) error {
	c, ok := m.Conn(p.Name())
	if !ok {
		return errors.New("could not find connection")
	}

	p.c = c
	return nil
}

func (p *Player) enableChunkCache() {
	if c, ok := p.Conn(); ok {
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
