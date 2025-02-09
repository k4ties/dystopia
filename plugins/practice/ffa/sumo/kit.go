package sumo

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var Kit = func() kit.Kit {
	var (
		items  = make(kit.Items)
		armour kit.Armour
	)

	items[0] = item.NewStack(item.Stick{}, 1).WithLore(text.Colourf("<red>dystopia</red>"))
	return kit.New(items, armour, effect.New(effect.Resistance, 255, handlers.SixThousandMinutes).WithoutParticles())
}()
