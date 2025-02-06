package instance

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"iter"
	"log/slog"
	"maps"
	"math"
	"slices"
	"sync"
)

type Impl struct {
	errorLog *slog.Logger

	players   map[uuid.UUID]*Player
	playersMu sync.RWMutex

	world    *world.World
	gameMode world.GameMode

	hidden     []hud.Element
	defaultRot cube.Rotation

	heightThresholdStatus atomic.Bool
	heightThreshold       int
	heightThresholdMode   OnIntersectThreshold
}

func (i *Impl) Messagef(s string, args ...any) {
	for p := range i.Players() {
		p.Message(text.Colourf(s, args...))
	}
}

func (i *Impl) HeightThresholdMode() OnIntersectThreshold {
	return i.heightThresholdMode
}

func (i *Impl) HeightThresholdEnabled() bool {
	return i.heightThresholdStatus.Load()
}

func (i *Impl) ToggleHeightThreshold() {
	i.heightThresholdStatus.Store(!i.heightThresholdStatus.Load())
}

func (i *Impl) HeightThreshold() int {
	return i.heightThreshold
}

func (i *Impl) Transfer(pl *Player, tx *world.Tx) {
	if i.Active(pl.UUID()) {
		//panic("cannot transfer player that is already active")
		return
	}

	pl.setTransferring(true)
	defer pl.setTransferring(false)

	FadeInCamera(pl.c, 0.5, false)

	if pl.Instance() != Nop {
		if imp, ok := pl.Instance().(interface {
			Hidden() []hud.Element
		}); ok {
			for _, elem := range imp.Hidden() {
				if !i.isHidden(elem) {
					_ = pl.ResetElements(elem)
				}
			}
		}

		pl.Instance().RemoveFromList(pl)
	}

	//if e, ok := pl.H().Entity(tx); !ok || tx == nil {
	//
	//}

	if !i.inWorld(tx) {
		h := tx.RemoveEntity(pl)

		i.World().Exec(func(tx *world.Tx) {
			tx.AddEntity(h)
		})
	}

	pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
		p.SetGameMode(i.gameMode)
		p.Teleport(i.World().Spawn().Vec3Centre())

		i.Rotate(p)
	})

	i.addToList(pl)
	_ = pl.HideElements(i.hidden...)
}

func (i *Impl) Rotate(p *player.Player) {
	currentRot := p.Rotation()
	maxRot := p.Rotation()

	yawDiff := findAngleDifference(currentRot.Yaw(), maxRot.Yaw()) - currentRot.Yaw()
	pitchDiff := findPitchDifference(currentRot.Pitch(), maxRot.Pitch()) - currentRot.Pitch()

	p.Move(mgl64.Vec3{}, yawDiff, pitchDiff)
}

func findPitchDifference(currentPitch float64, targetPitch float64) float64 {
	diff := targetPitch - currentPitch
	if math.Abs(diff) > 180 {
		if diff < 0 {
			diff += 360
		} else {
			diff -= 360
		}
	}
	return diff
}

func findAngleDifference(yaw float64, expectedYaw float64) float64 {
	diff := math.Mod(expectedYaw-yaw+1080, 360) - 180
	return diff
}

func (i *Impl) Hidden() []hud.Element {
	return i.hidden
}

func (i *Impl) addToList(p *Player) {
	if p == nil {
		return
	}

	if i.Active(p.UUID()) {
		panic("player is already in instance")
	}

	i.playersMu.Lock()
	i.players[p.UUID()] = p
	i.playersMu.Unlock()

	p.setInstance(i)
}

func (i *Impl) RemoveFromList(p *Player) {
	if !i.Active(p.UUID()) {
		panic("cannot remove from instance player that is not in instance")
	}

	i.playersMu.Lock()
	delete(i.players, p.UUID())
	i.playersMu.Unlock()

	p.setInstance(nil)
}

func (i *Impl) inWorld(tx *world.Tx) (found bool) {
	if tx.World() == i.World() {
		found = true
	}

	return
}

func (i *Impl) isHidden(e hud.Element) bool {
	return slices.Contains(i.hidden, e)
}

func (i *Impl) DefaultRotation() cube.Rotation {
	return i.defaultRot
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

	pl.enableChunkCache()
	return pl
}

type OnIntersectThreshold int

const (
	EventDeath OnIntersectThreshold = iota
	EventTeleportToSpawn
)

type HeightThresholdConfig struct {
	Enabled   bool
	Threshold int
	OnDeath   OnIntersectThreshold
}

var DisabledHeightThreshold = HeightThresholdConfig{
	Enabled: false,
}

func New(w *world.World, g world.GameMode, errorLogger *slog.Logger, defaultRot cube.Rotation, htc HeightThresholdConfig, hidden ...hud.Element) Instance {
	i := &Impl{players: make(map[uuid.UUID]*Player), world: w, gameMode: g, errorLog: errorLogger, defaultRot: defaultRot, hidden: hidden}

	if htc.Enabled {
		i.ToggleHeightThreshold()
		i.heightThreshold = htc.Threshold
		i.heightThresholdMode = htc.OnDeath
	}

	return i
}
