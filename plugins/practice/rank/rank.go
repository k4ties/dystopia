package rank

import (
	"github.com/google/uuid"
	"maps"
	"slices"
	"sync"
)

type Rank interface {
	Name() string
	Format() string

	DisplayRankName() bool
	Priority() Priority

	UUID() uuid.UUID
}

var ranks = struct {
	ranks map[string]Rank
	mu    sync.RWMutex
}{
	ranks: make(map[string]Rank),
}

func Register(r Rank) {
	ranks.mu.Lock()
	defer ranks.mu.Unlock()

	ranks.ranks[r.Name()] = r
}

func UnRegister(r Rank) {
	ranks.mu.Lock()
	defer ranks.mu.Unlock()

	delete(ranks.ranks, r.Name())
}

func ByName(name string) (Rank, bool) {
	ranks.mu.RLock()
	defer ranks.mu.RUnlock()

	r, ok := ranks.ranks[name]
	return r, ok
}

func MustByName(name string) Rank {
	r, ok := ByName(name)
	if !ok {
		return Player
	}

	return r
}

func ByUUID(uuid uuid.UUID) (Rank, bool) {
	ranks.mu.RLock()
	defer ranks.mu.RUnlock()

	for _, rank := range List() {
		if rank.UUID() == uuid {
			return rank, true
		}
	}

	return nil, false
}

func MustByUUID(uuid uuid.UUID) Rank {
	rank, ok := ByUUID(uuid)
	if !ok {
		return Player
	}

	return rank
}

func List() []Rank {
	ranks.mu.RLock()
	defer ranks.mu.RUnlock()

	return slices.Collect(maps.Values(ranks.ranks))
}
