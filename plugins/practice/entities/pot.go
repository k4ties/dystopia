package entities

import (
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"image/color"

	_ "unsafe"
)

func NewHealPotion(opts world.EntitySpawnOpts, owner world.Entity, colour color.RGBA) *world.EntityHandle {
	t := potion.StrongHealing()

	conf := splashPotionConf
	conf.Potion = t
	conf.Particle = particle.Splash{Colour: colour}
	conf.Hit = potionSplash(2.5, t, false)
	conf.Owner = owner.H()

	return opts.New(entity.SplashPotionType, conf)
}

var splashPotionConf = entity.ProjectileBehaviourConfig{
	Gravity: 0.045,
	Drag:    0.005,
	Damage:  -1,
	Sound:   sound.GlassBreak{},
}

//go:linkname potionSplash github.com/df-mc/dragonfly/server/entity.potionSplash
func potionSplash(float64, potion.Potion, bool) func(*entity.Ent, *world.Tx, trace.Result)
