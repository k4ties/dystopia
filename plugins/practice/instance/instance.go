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
	"strings"
	"sync"
)

type Instance interface {
	World() *world.World
	GameMode() world.GameMode

	ErrorLog() *slog.Logger

	NewPlayer(*player.Player) *Player
	Players() iter.Seq[*Player]
	Player(uuid.UUID) (*Player, bool)

	Active(uuid.UUID) bool
	Transfer(*Player, *world.Tx)

	addToList(*Player)
	removeFromList(*Player)

	HeightThresholdEnabled() bool
	HeightThresholdMode() OnIntersectThreshold

	ToggleHeightThreshold()
	HeightThreshold() int

	Messagef(string, ...any)
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
	instances.v[strings.ToLower(name)] = instance
}

func UnRegister(name string) {
	instances.mu.Lock()
	defer instances.mu.Unlock()
	delete(instances.v, strings.ToLower(name))
}

func ByName(name string) (Instance, bool) {
	instances.mu.RLock()
	defer instances.mu.RUnlock()

	instance, ok := instances.v[strings.ToLower(name)]
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
	pl.setInstance(nil)

	if err := pl.syncConn(plugin.M()); err != nil {
		Kick(p, ErrorSponge)
		return nil
	}

	pl.enableChunkCache()
	return pl
}

func LookupPlayer(pl *player.Player) (*Player, Instance) {
	for _, inst := range AllInstances() {
		if inst.Active(pl.UUID()) {
			for p := range inst.Players() {
				if p.UUID() == pl.UUID() {
					return p, inst
				}
			}
			break
		}
	}

	return nil, nil
}

func AllInstances() []Instance {
	instances.mu.RLock()
	defer instances.mu.RUnlock()

	return slices.Collect(maps.Values(instances.v))
}

func AllInstancesName() []string {
	var names []string

	instances.mu.RLock()
	for name := range instances.v {
		names = append(names, name)
	}
	instances.mu.RUnlock()

	return names
}

var Nop = nopInstance{}

var _ = (Instance)(Nop)

type nopInstance struct{}

func (n nopInstance) Player(uuid.UUID) (*Player, bool)          { return nil, false }
func (n nopInstance) Messagef(string, ...any)                   {}
func (n nopInstance) HeightThresholdMode() OnIntersectThreshold { return -1 }
func (n nopInstance) HeightThresholdEnabled() bool              { return false }
func (n nopInstance) ToggleHeightThreshold()                    {}
func (n nopInstance) HeightThreshold() int                      { return -1 }
func (n nopInstance) World() *world.World                       { return nil }
func (n nopInstance) GameMode() world.GameMode                  { return nil }
func (n nopInstance) ErrorLog() *slog.Logger                    { return nil }
func (n nopInstance) NewPlayer(*player.Player) *Player          { return nil }
func (n nopInstance) Players() iter.Seq[*Player]                { return nil }
func (n nopInstance) Active(uuid.UUID) bool                     { return false }
func (n nopInstance) Transfer(p *Player, tx *world.Tx) {
	p.Instance().removeFromList(p)
}
func (n nopInstance) addToList(*Player)      {}
func (n nopInstance) removeFromList(*Player) {}
