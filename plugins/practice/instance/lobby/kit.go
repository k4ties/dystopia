package lobby

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/k4ties/dystopia/plugins/practice/kit"
)

var Kit = func() kit.Kit {
	stack := item.NewStack(item.Sword{Tier: item.ToolTierNetherite}, 1)

	var items = make(kit.Items)
	items[4] = kit.ApplyIdentifier(kit.FFAIdentifier, kit.FillNames("FFA", stack))

	return kit.New(items, kit.NopArmour)
}()
