package handlers

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"time"
)

const SixThousandMinutes = time.Hour * 10000

type Practice struct {
	Decliner
	plugin.NopPlayerHandler
}

func (Practice) HandleSpawn(p *player.Player) {
	p.SetImmobile()

	if c, ok := plugin.M().Conn(p.Name()); ok {
		instance.FadeInCamera(c, 1.5, false)
	}

	resetPlayer(p)
	sendDefaultEffects(p)

	pl := instance.NewPlayer(p)
	lobby.Instance().Transfer(pl, p.Tx())
}

func (Practice) HandleQuit(p *player.Player) {
	if pl := instance.LookupPlayer(p); pl != nil {
		if i := pl.Instance(); i != instance.Nop {
			i.RemoveFromList(pl)
		}
	}
}

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

func resetPlayer(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()

	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}

	p.SetHeldItems(item.Stack{}, item.Stack{}) // clears totem
	p.SetGameMode(world.GameModeSurvival)

	p.CloseForm()
	p.CloseDialogue()

	p.EnableInstantRespawn()
	p.SetMobile()

	p.ShowCoordinates()
}

func welcomePlayer(p *player.Player) {
	_ = p.SetHeldSlot(4)

}

//future fps counter implementation
//p.UpdateDiagnostics(session.Diagnostics{
//AverageFramesPerSecond:        -1,
//AverageServerSimTickTime:      -1,
//AverageClientSimTickTime:      -1,
//AverageBeginFrameTime:         -1,
//AverageInputTime:              -1,
//AverageRenderTime:             -1,
//AverageEndFrameTime:           -1,
//AverageRemainderTimePercent:   -1,
//AverageUnaccountedTimePercent: -1,
//})
