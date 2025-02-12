package gapple

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var Kit = func() kit.Kit {
	var (
		items  = make(kit.Items)
		armour kit.Armour
	)

	items[0] = item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Sharpness, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	items[1] = kit.ApplyIdentifier(kit.GAppleIdentifier, item.NewStack(item.GoldenApple{}, 16).WithLore(text.Colourf("<red>dystopia</red>")))

	armour[0] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[1] = item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[2] = item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[3] = item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))

	return kit.New(items, armour)
}()
