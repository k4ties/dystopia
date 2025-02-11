package items

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/k4ties/dystopia/plugins/practice/entities"
	"image/color"
)

func NewHealingPotion() world.Item {
	return HealingPotion{i: item.Potion{Type: potion.StrongHealing()}}
}

type HealingPotion struct {
	i item.Potion
}

func (s HealingPotion) MaxCount() int {
	return 1
}

func (s HealingPotion) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	create := entities.NewHealPotion
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(0.5)}

	tx.AddEntity(create(opts, user, color.RGBA{A: 255}))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}

// EncodeItem ...
func (s HealingPotion) EncodeItem() (name string, meta int16) {
	return "minecraft:splash_potion", int16(s.i.Type.Uint8())
}
