package entities

import (
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

func NewEnderPearl(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := enderPearlConf
	conf.Owner = owner.H()
	return opts.New(entity.EnderPearlType, conf)
}

var enderPearlConf = entity.ProjectileBehaviourConfig{
	Gravity: 0.03,
	Drag:    0.01,
	// Particle: particle.EndermanTeleport{},
	Sound: sound.Teleport{},
	Hit:   teleport,
}

func teleport(e *entity.Ent, tx *world.Tx, target trace.Result) {
	owner, _ := e.Behaviour().(*entity.ProjectileBehaviour).Owner().Entity(tx)
	if user, ok := owner.(interface {
		Teleport(pos mgl64.Vec3)
		entity.Living
	}); ok {
		tx.PlaySound(user.Position(), sound.Teleport{})
		user.Teleport(target.Position())
		// user.Hurt(5, FallDamageSource{})
	}
}
