package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

var kits = struct {
	v  map[string]Kit
	mu sync.RWMutex
}{
	v: make(map[string]Kit),
}

func Register(kit Kit, name string) {
	kits.mu.Lock()
	defer kits.mu.Unlock()
	kits.v[name] = kit
}

func ByName(name string) (Kit, bool) {
	kits.mu.RLock()
	defer kits.mu.RUnlock()

	v, ok := kits.v[name]
	return v, ok
}

func UnRegister(name string) {
	kits.mu.Lock()
	defer kits.mu.Unlock()

	delete(kits.v, name)
}

type Armour [4]item.Stack
type Items map[int]ItemEntry

func NewArmour(helmet, chestplate, leggings, boots item.Stack) Armour {
	return Armour{
		helmet,
		chestplate,
		leggings,
		boots,
	}
}

type ItemEntry struct {
	Slot int
	Item item.Stack
}

func NewItemEntry(slot int, item item.Stack) ItemEntry {
	return ItemEntry{slot, item}
}

func NewItems(items ...ItemEntry) Items {
	if len(items) > 36 {
		panic("cant have more than 36 items")
	}

	var m Items

	for _, entry := range items {
		m[entry.Slot] = entry
	}

	return m
}

type Kit interface {
	Items() Items
	Armour() Armour
	Effects() []effect.Effect
}

func Send(k Kit, p *player.Player) {
	SendItems(k, p)
	SendArmour(k, p)
	SendEffects(k, p)
}

func SendItems(k Kit, p *player.Player) {
	for slot, i := range k.Items() {
		_ = p.Inventory().SetItem(slot, i.Item)
	}
}

func SendArmour(k Kit, p *player.Player) {
	armour := p.Armour()

	armour.SetHelmet(k.Armour()[0])
	armour.SetChestplate(k.Armour()[1])
	armour.SetLeggings(k.Armour()[2])
	armour.SetBoots(k.Armour()[3])
}

func SendEffects(k Kit, p *player.Player) {
	for _, e := range k.Effects() {
		p.AddEffect(e)
	}
}
