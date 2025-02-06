package whitelist

import (
	"slices"
	"sync"
	"sync/atomic"
)

var w = struct {
	enabled atomic.Bool

	players   []string
	playersMu sync.RWMutex
}{}

func Toggle() {
	w.enabled.Store(!w.enabled.Load())
}

func Setup(whitelisted bool, players ...string) {
	w.enabled.Store(whitelisted)
	w.players = append(w.players, players...)
}

func Add(player string) {
	w.playersMu.Lock()
	defer w.playersMu.Unlock()

	w.players = append(w.players, player)
}

func Enabled() bool {
	return w.enabled.Load()
}

func Whitelisted(player string) bool {
	if w.enabled.Load() {
		w.playersMu.RLock()
		defer w.playersMu.RUnlock()
		return slices.Contains(w.players, player)
	}

	return false
}

func Remove(player string) {
	w.playersMu.Lock()
	defer w.playersMu.Unlock()

	w.players = remove(w.players, slices.Index(w.players, player))
}

func remove(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
