package instance

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"iter"
	"log/slog"
	"maps"
	"slices"
	"sync"
)

type Instance interface {
	World() *world.World
	GameMode() world.GameMode

	NewPlayer(*player.Player) *Player

	Players() iter.Seq[*Player]
	Active(uuid.UUID) bool

	AddPlayer(*Player)
	RemovePlayer(*Player)

	ErrorLog() *slog.Logger
}

var instances = struct {
	v  map[string]Instance
	mu sync.RWMutex
}{
	v: make(map[string]Instance),
}

func Register(name string, instance Instance) {
	instances.mu.Lock()
	defer instances.mu.Unlock()
	instances.v[name] = instance
}

func UnRegister(name string) {
	instances.mu.Lock()
	defer instances.mu.Unlock()
	delete(instances.v, name)
}

func ByName(name string) (Instance, bool) {
	instances.mu.RLock()
	defer instances.mu.RUnlock()

	instance, ok := instances.v[name]
	return instance, ok
}

func MustByName(name string) Instance {
	i, ok := ByName(name)
	if !ok {
		panic("instance not found: " + name)
	}

	return i
}

func NewPlayer(p *player.Player) *Player {
	pl := &Player{Player: p}

	if err := pl.syncConn(plugin.M()); err != nil {
		Kick(p, ErrorSponge)
		return nil
	}

	return pl
}

func LookupPlayer(pl *player.Player) *Player {
	allInstances := slices.Collect(maps.Values(instances.v))

	for _, inst := range allInstances {
		if inst.Active(pl.UUID()) {
			for p := range inst.Players() {
				if p.UUID() == pl.UUID() {
					return p
				}
			}
			break
		}
	}

	return nil
}
