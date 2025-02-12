package handlers

import (
	_ "embed"
	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/dystopia/embeddable"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/items"
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"regexp"
	"strings"
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

	lastChattedAt atomic.Value[time.Time]
	lastCommandAt atomic.Value[time.Time]

	lastMessage atomic.Value[string]
}

const (
	JoinFormat = "<dark-grey>(<green>+</green>)</dark-grey> %s"
	QuitFormat = "<dark-grey>(<red>-</red>)</dark-grey> %s"
)

func (pr *Practice) HandleSpawn(p *player.Player) {
	_, _ = chat.Global.WriteString(text.Colourf(JoinFormat, p.Name()))
	p.SetNameTag(text.Grey + p.Name())

	if c, ok := plugin.M().Conn(p.Name()); ok {
		instance.FadeInCamera(c, 1.5, false)
	}

	welcomePlayer(p)
	_ = pr.spawnRoutine(p, p.Tx())
}

func (pr *Practice) spawnRoutine(p *player.Player, tx *world.Tx) *instance.Player {
	// player must be in lobby instance on spawn
	pl := instance.NewPlayer(p)

	pr.i.Transfer(pl, tx)
	if pr.i.Active(p.UUID()) {
		p.Teleport(pr.i.World().Spawn().Vec3Centre())
	}

	pl.SendKit(lobby.Kit, tx)
	return pl
}

func (pr *Practice) HandleMove(ctx *player.Context, newPos mgl64.Vec3, _ cube.Rotation) {
	if p, i := instance.LookupPlayer(ctx.Val()); p != nil && i != nil {
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

func (pr *Practice) HandleQuit(p *player.Player) {
	_, _ = chat.Global.WriteString(text.Colourf(QuitFormat, p.Name()))

	if fpl, in := ffa.LookupPlayer(p); fpl != nil && in != nil {
		in.ResetPearlCooldown(fpl)
	}

	if pl, in := instance.LookupPlayer(p); pl != nil && in != nil {
		if i := pl.Instance(); i != instance.Nop {
			instance.Nop.Transfer(pl, nil)
		}
	}
}

func countPots(inv *inventory.Inventory) int {
	var count int

	for _, i := range inv.Items() {
		if _, isPot := i.Item().(items.HealingPotion); isPot {
			count++
		}
	}

	return count
}

func countGApples(inv *inventory.Inventory) int {
	var count int

	for _, i := range inv.Items() {
		if _, isGApple := i.Item().(item.GoldenApple); isGApple {
			count += i.Count()
		}
	}

	return count
}

func (pr *Practice) HandleDeath(deadPlayer *player.Player, src world.DamageSource, keepInv *bool) {
	*keepInv = true
	deadPlayerPl, deadPlayerIn := ffa.LookupPlayer(deadPlayer) // player can only be dead on ffa

	deadPlayerGapplesBeforeDeath := countGApples(deadPlayer.Inventory())
	deadPlayerPotsBeforeDeath := countPots(deadPlayer.Inventory())

	if deadPlayerPl == nil || deadPlayerIn == nil {
		return
	}

	if _, ok := src.(intersectThresholdCause); ok {
		deadPlayerIn.Messagef("<red>%s</red> fell into the void", deadPlayer.Name())
	}

	var killerName = "..."
	var killerIsInstancePlayer = false

	var (
		killerPl *instance.Player
		killerIn *ffa.Instance
	)

	if a, ok := src.(entity.AttackDamageSource); ok {
		if attackerP, ok := a.Attacker.(*player.Player); ok {
			pl, in := ffa.LookupPlayer(attackerP)
			if pl != nil && in != nil {
				killerIsInstancePlayer = true
				killerName = pl.Name()

				killerPl = pl
				killerIn = in
			}
		}
	}

	deadPlayerPl.Reset(nil, func(p *player.Player) {
		p.SetGameMode(world.GameModeSpectator)
	}, func(p *player.Player) {
		t := title.New(text.Red + "YOU'RE DEAD")
		t = t.WithSubtitle(text.Grey + killerName)
		t = t.WithDuration(time.Second * 3)
		p.SendTitle(t)
	})

	if killerIsInstancePlayer {
		var msg = formatDeadMessage(DeadMessageBlank, deadPlayer.Name(), killerName)

		if kitIncludesPotions(deadPlayerIn.Kit()) {
			msg = formatDeadMessage(DeadMessageCount, deadPlayer.Name(), deadPlayerPotsBeforeDeath, "POTS", killerName, countPots(killerPl.Inventory()), "POTS")
		} else if kitIncludesGApples(deadPlayerIn.Kit()) {
			msg = formatDeadMessage(DeadMessageCount, deadPlayer.Name(), deadPlayerGapplesBeforeDeath, "GAPPLES", killerName, countGApples(killerPl.Inventory()), "GAPPLES")
		}

		deadPlayerIn.Messagef(msg)
		killerIn.ReKit(killerPl, nil)
	}

	time.AfterFunc(time.Second*4, func() {
		if plugin.M().Online(deadPlayer.UUID()) {
			deadPlayerPl.ExecSafe(func(p *player.Player, tx *world.Tx) {
				lobby.TransferWithRoutine(deadPlayerPl, tx)
			})
		}
	})
}

const (
	DeadMessageBlankFormat = "<red>%s</red> was killed by <red>%s</red>"
	DeadMessageCountFormat = "<red>%s</red> [%d %s] was killed by <red>%s</red> [%d %s]"
)

const (
	DeadMessageBlank = iota
	DeadMessageCount
)

func formatDeadMessage(index int, args ...any) string {
	switch index {
	case DeadMessageBlank:
		if len(args) < 2 {
			panic("not enough arguments")
		}
		return text.Colourf(DeadMessageBlankFormat, args...)
	case DeadMessageCount:
		if len(args) < 6 {
			panic("not enough arguments")
		}
		return text.Colourf(DeadMessageCountFormat, args...)
	}

	panic("should never happen")
}

func kitIncludesPotions(k kit.Kit) bool {
	for _, i := range k.Items() {
		if _, isPot := i.Item().(items.HealingPotion); isPot {
			return true
		}
	}

	return false
}

func kitIncludesGApples(k kit.Kit) bool {
	for _, i := range k.Items() {
		if _, isGApple := i.Item().(item.GoldenApple); isGApple {
			return true
		}
	}

	return false
}

func (pr *Practice) HandleItemUseOnBlock(ctx *player.Context, _ cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	reHandleItemUseOnBlock(ctx)
}

func (pr *Practice) HandleItemUse(ctx *player.Context) {
	reHandleItemUse(ctx)
}

func reHandleItemUseOnBlock(ctx *player.Context) {
	i, _ := ctx.Val().HeldItems()

	switch kit.LoadIdentifier(i) {
	case kit.GAppleIdentifier, kit.PotIdentifier, kit.PearlIdentifier:
		// do nothing, just allow itemUse to re handle it
	default:
		ctx.Cancel()
	}
}

func reHandleItemUse(ctx *player.Context) {
	i, _ := ctx.Val().HeldItems()

	switch kit.LoadIdentifier(i) {
	case kit.FFAIdentifier:
		ctx.Val().SendForm(ffa.NewForm())
		ctx.Cancel()
	case kit.NodebuffSwordIdentifier:
		ctx.Cancel()
	case kit.PearlIdentifier:
		ctx.Val().SetCooldown(items.Pearl{}, -1)

		if pl, in := ffa.LookupPlayer(ctx.Val()); pl != nil && in != nil {
			if in.HasPearCooldown() {
				in.MustCoolDown(pl.UUID(), func(c *cooldown.CoolDown) {
					if c.Active() {
						ctx.Val().SendJukeboxPopup(text.Reset + text.Colourf("<red>Please wait %d more seconds to use pearl again", int(c.Remaining().Seconds())))
						ctx.Cancel()
						return
					}

					in.StartPearlCoolDown(pl, nil)
					ctx.Val().Messagef(text.Colourf("<red>Pearl cooldown has started</red>"))
				})
			}
		}
	}
}

func (pr *Practice) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if pl, in := ffa.LookupPlayer(ctx.Val()); pl != nil && in != nil {
		if _, ok := src.(effect.InstantHealingSource); ok {
			*health *= 2.5
		}
	}
}

func (pr *Practice) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	if time.Since(pr.lastChattedAt.Load()) <= time.Second/2 {
		ctx.Val().Messagef(text.Colourf("<red>Please don't spam</red>"))
		return
	}

	if len(*msg) > 130 {
		ctx.Val().Messagef(text.Colourf("<red>You've ran out of characters. (%v/130)</red>", len(*msg)))
		return
	}

	if notAlphaOnly(*msg) {
		return
	}

	*msg = removeExtraSpaces(*msg)
	*msg = alphaReplacer(*msg)
	*msg = strings.TrimSpace(*msg)

	if pr.lastMessage.Load() == *msg {
		ctx.Val().Messagef(text.Colourf("<red>Please don't send identical messages</red>"))
		return
	}

	pr.lastChattedAt.Store(time.Now())

	if *msg != "" {
		_, _ = chat.Global.WriteString(text.Colourf("<grey>%s:</grey> %s", ctx.Val().Name(), *msg))
		pr.lastMessage.Store(*msg)
	}
}

func notAlphaOnly(s string) bool {
	return len(alphaReplacer(s)) == 0
}

func alphaReplacer(s string) string {
	re := regexp.MustCompile(`[\s§<>]+`)
	return re.ReplaceAllString(s, " ")
}

func removeExtraSpaces(s string) string {
	words := strings.Fields(s)
	return strings.Join(words, " ")
}

func (pr *Practice) HandleCommandExecution(ctx *player.Context, _ cmd.Command, _ []string) {
	if time.Since(pr.lastCommandAt.Load()) <= time.Second/2 {
		ctx.Val().Messagef(text.Colourf("<red>Please don't spam</red>"))
		ctx.Cancel()
		return
	}

	pr.lastCommandAt.Store(time.Now())
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

func (pr *Practice) HandleHurt(ctx *player.Context, damage *float64, _ bool, attackImmunity *time.Duration, src world.DamageSource) {
	attackSrc, isAttackSource := src.(entity.AttackDamageSource)
	_, isProjectileSource := src.(entity.ProjectileDamageSource)

	attacker, isPlayer := attackSrc.Attacker.(*player.Player)

	if lobby.Instance().Active(ctx.Val().UUID()) {
		ctx.Cancel()
		return
	}

	if isPlayer && attacker.GameMode() == world.GameModeSpectator {
		ctx.Cancel()
		return
	}

	if !isAttackSource && !isProjectileSource {
		ctx.Cancel()
		return
	}

	*attackImmunity = time.Duration(KnockbackConfig.KnockBack.Immunity) * time.Millisecond

	if isProjectileSource {
		*attackImmunity = 0
	}

	if isPlayer && isCritical(attacker) {
		reHandleCritical(attacker, ctx.Val(), ctx.Val().Tx())
		*damage *= 1.5
	}
}

const scoreTagFormat = "\uE10C %d <red>|</red> %dms"

func scoreTagTask(i instance.Instance, p *instance.Player) {
	if i == lobby.Instance() || p.GameMode() == world.GameModeSpectator {
		p.ExecSafe(func(p *player.Player, tx *world.Tx) {
			p.SetScoreTag("")
		})
		return
	}

	p.ExecSafe(func(p *player.Player, tx *world.Tx) {
		p.SetScoreTag(text.Colourf(scoreTagFormat, int(p.Health()), p.Latency().Milliseconds()))
	})
}

func (pr *Practice) HandleAttackEntity(_ *player.Context, _ world.Entity, force, height *float64, critical *bool) {
	*force = KnockbackConfig.KnockBack.Force
	*height = KnockbackConfig.KnockBack.Height

	*critical = false // prevents double hit 🤡
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

func welcomePlayer(p *player.Player) {
	_ = p.SetHeldSlot(4)

	welcomeTitle := title.New(text.Colourf("<red>Dystopia</red>"))
	welcomeTitle = welcomeTitle.WithFadeInDuration(time.Second * 2).WithDuration(time.Second).WithFadeOutDuration(time.Second)
	welcomeTitle = welcomeTitle.WithSubtitle(text.Colourf("<white>Welcome, <grey>%s</grey>!</white>", p.Name()))

	p.SendTitle(welcomeTitle)
}
