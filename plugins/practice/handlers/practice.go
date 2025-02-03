package handlers

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"time"
)

const SixThousandMinutes = time.Hour * 10000

type Practice struct {
	Decliner
	plugin.NopPlayerHandler
}

func (Practice) HandleSpawn(p *player.Player) {
	resetPlayer(p)
	sendDefaultEffects(p)
}

func (Practice) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if lobby.Instance().Active(ctx.Val().UUID()) {
		ctx.Cancel()
	}
}

func sendDefaultEffects(p *player.Player) {
	p.AddEffect(effect.New(effect.NightVision, 1, SixThousandMinutes).WithoutParticles())
}

func resetPlayer(p *player.Player) {
	p.SetMobile()

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
