package instance

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"iter"
	"log/slog"
	"maps"
	"sync"
)

type Impl struct {
	errorLog *slog.Logger

	players   map[uuid.UUID]*Player
	playersMu sync.RWMutex

	world    *world.World
	gameMode world.GameMode
}

func (i *Impl) World() *world.World {
	return i.world
}

func (i *Impl) GameMode() world.GameMode {
	return i.gameMode
}

func (i *Impl) Players() iter.Seq[*Player] {
	i.playersMu.RLock()
	defer i.playersMu.RUnlock()
	return maps.Values(i.players)
}

func (i *Impl) AddPlayer(pl *Player) {
	if pl == nil {
		return
	}

	if i.Active(pl.UUID()) {
		panic("player is already in instance")
	}

	if pl.Instance() != nil {
		pl.Instance().RemovePlayer(pl)
	}

	pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
		tx.RemoveEntity(p) // remove from current world

		p.SetGameMode(i.gameMode)
		p.Teleport(i.World().Spawn().Vec3Centre())

		i.World().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
	})

	i.playersMu.Lock()
	i.players[pl.UUID()] = pl
	i.playersMu.Unlock()

	pl.setInstance(i)
}

func (i *Impl) RemovePlayer(pl *Player) {
	if !i.Active(pl.UUID()) {
		panic("player is not in instance")
	}

	pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
		tx.RemoveEntity(p)
	})

	i.playersMu.Lock()
	delete(i.players, pl.UUID())
	i.playersMu.Unlock()

	pl.setInstance(nil)
}

func (i *Impl) Active(u uuid.UUID) bool {
	i.playersMu.RLock()
	defer i.playersMu.RUnlock()

	_, ok := i.players[u]
	return ok
}

func (i *Impl) ErrorLog() *slog.Logger {
	return i.errorLog
}

func (i *Impl) NewPlayer(p *player.Player) *Player {
	pl := &Player{Player: p, instance: i}

	if err := pl.syncConn(plugin.M()); err != nil {
		Kick(p, ErrorSponge)
		return nil
	}

	return pl
}

func New(w *world.World, g world.GameMode, errorLogger *slog.Logger) Instance {
	return &Impl{players: make(map[uuid.UUID]*Player), world: w, gameMode: g, errorLog: errorLogger}
}
