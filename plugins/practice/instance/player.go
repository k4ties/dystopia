package instance

import (
	"errors"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
	"sync"
)

type Player struct {
	*player.Player

	instance   Instance
	instanceMu sync.Mutex

	c session.Conn
}

func (p *Player) setInstance(i Instance) {
	p.instanceMu.Lock()
	defer p.instanceMu.Unlock()
	p.instance = i
}

func (p *Player) Instance() Instance {
	p.instanceMu.Lock()
	defer p.instanceMu.Unlock()
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
