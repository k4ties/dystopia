package nodebuff

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var Kit = func() kit.Kit {
	var (
		items   = make(kit.Items)
		armour  kit.Armour
		effects []effect.Effect
	)

	for i := 1; i <= 36; i++ {
		var added item.Stack

		switch i {
		case 1:
			added = item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Sharpness, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
		case 2:
			added = item.NewStack(item.EnderPearl{}, 16).WithLore(text.Colourf("<red>dystopia</red>"))
		default:
			added = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1).WithLore(text.Colourf("<red>dystopia</red>"))
		}

		items[i-1] = added
	}

	armour[0] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[1] = item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[2] = item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))
	armour[3] = item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 4)).WithLore(text.Colourf("<red>dystopia</red>"))

	effects = append(effects, effect.New(effect.Speed, 1, handlers.SixThousandMinutes).WithoutParticles())
	return kit.New(items, armour, effects...)
}()
