package handlers

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/dystopia/embeddable"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

const SixThousandMinutes = time.Hour * 10000

func NewPractice(i instance.Instance) *Practice {
	return &Practice{i: i}
}

type Practice struct {
	Decliner
	plugin.NopPlayerHandler

	i instance.Instance // default instance
}

func (pr Practice) HandleSpawn(p *player.Player) {
	if c, ok := plugin.M().Conn(p.Name()); ok {
		instance.FadeInCamera(c, 1.5, false)
	}

	pr.spawnRoutine(p, p.Tx())
}

func (pr Practice) spawnRoutine(p *player.Player, tx *world.Tx) {
	// player must be in lobby instance on spawn

	pl := instance.NewPlayer(p)

	pr.i.Transfer(pl, tx)
	if pr.i.Active(p.UUID()) {
		p.Teleport(pr.i.World().Spawn().Vec3Centre())
	}

	pl.SendKit(lobby.Kit, tx)
}

func (pr Practice) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	if p := instance.LookupPlayer(ctx.Val()); p != nil {
		if i := p.Instance(); i != instance.Nop {
			if i.HeightThresholdEnabled() {
				if newPos.Y() <= float64(i.HeightThreshold()) && p.GameMode() == i.GameMode() && !p.Transferring() {
					switch i.HeightThresholdMode() {
					case instance.EventTeleportToSpawn:
						i.Transfer(p, p.Tx())
						pr.spawnRoutine(ctx.Val(), nil)
					case instance.EventDeath:
						var nop bool
						pr.HandleDeath(ctx.Val(), intersectThresholdCause{}, &nop)
					}
				}
			}
		}
	}
}

type intersectThresholdCause struct{}

func (i intersectThresholdCause) ReducedByArmour() bool     { return false }
func (i intersectThresholdCause) ReducedByResistance() bool { return false }
func (i intersectThresholdCause) Fire() bool                { return false }

func (Practice) HandleQuit(p *player.Player) {
	if pl := instance.LookupPlayer(p); pl != nil {
		if i := pl.Instance(); i != instance.Nop {
			i.RemoveFromList(pl)
		}
	}
}

func (pr Practice) HandleDeath(dfp *player.Player, src world.DamageSource, keepInv *bool) {
	*keepInv = true
	dfp.Respawn()

	if pl := instance.LookupPlayer(dfp); pl != nil {
		pl.SendKit(kit.Empty, dfp.Tx())
	}

	if _, ok := src.(intersectThresholdCause); ok {
		if pl := instance.LookupPlayer(dfp); pl != nil {
			if i := pl.Instance(); i != instance.Nop {
				i.Messagef("<red>%s</red> fell into the void.", dfp.Name())
			}
		}
	}

	var killerName = "..."
	if a, ok := src.(entity.AttackDamageSource); ok {
		if p, ok := a.Attacker.(*player.Player); ok {
			killerName = p.Name()
		}
	}

	if killerName != "..." {
		if pl := instance.LookupPlayer(dfp); pl != nil {
			if i := pl.Instance(); i != instance.Nop {
				i.Messagef("<red>%s</red> was killed by <red>%s</red>", dfp.Name(), killerName)
			}
		}
	}

	dur := time.Millisecond * 500

	deadTitle := title.New(text.Colourf("<red>YOU ARE DEAD</red>"))
	deadTitle = deadTitle.WithDuration(dur * 5).WithFadeInDuration(dur).WithFadeOutDuration(dur).WithSubtitle(killerName)

	dfp.SendTitle(deadTitle)
	dfp.SetGameMode(world.GameModeSpectator)

	dfp.Teleport(dfp.Position().Add(mgl64.Vec3{3, 2, 3}))

	dfp.SetInvisible()
	dfp.SetScale(0)

	time.AfterFunc(time.Second*3, func() {
		if plugin.M().Online(dfp.UUID()) {
			if pl := instance.LookupPlayer(dfp); pl != nil {
				pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
					pr.i.Transfer(pl, tx)
					pr.spawnRoutine(p, tx)
				})
			}
		}
	})
}

func (Practice) HandleItemUse(ctx *player.Context) {
	i, _ := ctx.Val().HeldItems()

	switch kit.LoadIdentifier(i) {
	case lobby.KitFFAItemIdentifier:
		ctx.Val().SendForm(ffa.NewForm())
	}
}

type kbConfig struct {
	KnockBack struct {
		Force    float64
		Height   float64
		Immunity int
	}
}

//go:embed knockback.json
var knockbackConfig []byte

var KnockbackConfig = embeddable.MustJSON[kbConfig](knockbackConfig)

func (Practice) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if lobby.Instance().Active(ctx.Val().UUID()) {
		ctx.Cancel()
		return
	}

	attackSrc, isAttackSource := src.(entity.AttackDamageSource)
	_, isProjectileSource := src.(entity.ProjectileDamageSource)

	if !isAttackSource && !isProjectileSource {
		ctx.Cancel()
		return
	}

	*attackImmunity = time.Duration(KnockbackConfig.KnockBack.Immunity) * time.Millisecond

	if isProjectileSource {
		*attackImmunity = 0
	}

	if p, ok := attackSrc.Attacker.(*player.Player); ok && isCritical(p) {
		*damage *= 1.5
	}
}

func (Practice) HandleAttackEntity(ctx *player.Context, attacked world.Entity, force, height *float64, crit *bool) {
	if lobby.Instance().Active(ctx.Val().UUID()) {
		ctx.Cancel()
		return
	}

	*force = KnockbackConfig.KnockBack.Force
	*height = KnockbackConfig.KnockBack.Height

	*crit = false // prevents double hit 🤡

	// since we've cancelled critical, we need to handle it by ourselves
	p := ctx.Val()

	if isCritical(p) {
		reHandleCritical(p, attacked, p.Tx())
	}
}

func reHandleCritical(p *player.Player, attacked world.Entity, tx *world.Tx) {
	for _, v := range tx.Viewers(p.Position()) {
		v.ViewEntityAction(attacked, entity.CriticalHitAction{})
	}
}

func isCritical(p *player.Player) bool {
	_, slowFalling := p.Effect(effect.SlowFalling)
	_, blind := p.Effect(effect.Blindness)

	return !p.Sprinting() && !p.Flying() && p.FallDistance() > 0 && !slowFalling && !blind
}

func sendDefaultEffects(p *player.Player) {
	p.AddEffect(effect.New(effect.NightVision, 1, SixThousandMinutes).WithoutParticles())
}

func welcomePlayer(p *player.Player) {
	_ = p.SetHeldSlot(4)

	welcomeTitle := title.New(text.Colourf("<red>Dystopia</red>"))
	welcomeTitle = welcomeTitle.WithFadeInDuration(time.Second * 2).WithDuration(time.Second).WithFadeOutDuration(time.Second)
	welcomeTitle = welcomeTitle.WithSubtitle(text.Colourf("Welcome, <grey>%s</grey>!", p.Name()))

	p.SendTitle(welcomeTitle)
}
