package handlers

import (
	_ "embed"
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
	"github.com/k4ties/dystopia/plugins/practice/kit"
	"github.com/k4ties/dystopia/plugins/practice/rank"
	"github.com/k4ties/dystopia/plugins/practice/user"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
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
	pl := pr.spawnRoutine(p, p.Tx())

	if c := pl.MustConn(); c != nil {
		if _, ok := pl.User(); !ok {
			user.Register(user.New(c, user.OfflineUser{
				IPS:       []string{p.Addr().String()},
				DeviceIDS: []string{p.DeviceID()},
				FirstJoin: time.Now(),
				Name:      p.Name(),
				XUID:      p.XUID(),
				UUID:      p.UUID().String(),
				Deaths:    0,
				Kills:     0,
				Rank:      rank.Player,
			}))
		}
	}
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

func (pr *Practice) HandleQuit(p *player.Player) {
	_, _ = chat.Global.WriteString(text.Colourf(QuitFormat, p.Name()))

	if pl := instance.LookupPlayer(p); pl != nil {
		if i := pl.Instance(); i != instance.Nop {
			i.RemoveFromList(pl)
		}

		if u, ok := pl.User(); ok {
			u.Dystopia().SyncQuit()
		}
	}
}

func countPots(inv *inventory.Inventory) int {
	var count int

	for _, i := range inv.Items() {
		if _, isPot := i.Item().(item.SplashPotion); isPot {
			count++
		}
	}

	return count
}

func (pr *Practice) HandleDeath(dfp *player.Player, src world.DamageSource, keepInv *bool) {
	*keepInv = true
	potsBeforeDeath := countPots(dfp.Inventory())

	pl := instance.LookupPlayer(dfp)
	if pl == nil {
		return
	}

	pl.SendKit(kit.Empty, dfp.Tx())

	if u, ok := pl.User(); ok {
		u.Dystopia().Death()
	}

	if _, ok := src.(intersectThresholdCause); ok {
		if i := pl.Instance(); i != instance.Nop {
			i.Messagef("<red>%s</red> fell into the void.", dfp.Name())
		}
	}

	var killerName = "..."
	var killer *player.Player

	var killerIn *instance.Player

	if a, ok := src.(entity.AttackDamageSource); ok {
		if p, ok := a.Attacker.(*player.Player); ok {
			killerName = p.Name()
			killer = p

			if pl := instance.LookupPlayer(p); pl != nil {
				killerIn = pl

				if u, ok := pl.User(); ok {
					u.Dystopia().Kill()
				}
			}
		}
	}

	if killerName != "..." {
		if i := pl.Instance(); i != instance.Nop {
			plr, in := ffa.LookupPlayer(dfp)

			if plr != nil || in != nil {
				if kitIncludesPotions(in.Kit()) {
					i.Messagef("<red>%s</red> [%d POTS] was killed by <red>%s</red> [%d POTS]", dfp.Name(), potsBeforeDeath, killerName, countPots(killer.Inventory()))
				} else {
					i.Messagef("<red>%s</red> was killed by <red>%s</red>", dfp.Name(), killerName)
				}

				if killerIn != nil {
					killerIn.SendKit(in.Kit(), killer.Tx())
				}
			}
		}
	}

	dur := time.Second / 2

	deadTitle := title.New(text.Colourf("<red>YOU ARE DEAD</red>"))
	deadTitle = deadTitle.WithDuration(dur * 5).WithFadeInDuration(dur).WithFadeOutDuration(dur).WithSubtitle(killerName)

	dfp.SendTitle(deadTitle)
	dfp.SetGameMode(world.GameModeSpectator)

	time.AfterFunc(time.Second*3, func() {
		if plugin.M().Online(dfp.UUID()) {
			pl.ExecSafe(func(p *player.Player, tx *world.Tx) {
				pr.i.Transfer(pl, tx)
				pr.spawnRoutine(p, tx)
			})
		}
	})
}

func kitIncludesPotions(k kit.Kit) bool {
	for _, i := range k.Items() {
		if _, isPot := i.Item().(item.SplashPotion); isPot {
			return true
		}
	}

	return false
}

func (pr *Practice) HandleItemUse(ctx *player.Context) {
	i, _ := ctx.Val().HeldItems()

	switch kit.LoadIdentifier(i) {
	case kit.FFAIdentifier:
		ctx.Val().SendForm(ffa.NewForm())
	case kit.PearlIdentifier:
		ctx.Val().Messagef("pearl used")
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

func (pr *Practice) HandleCommandExecution(ctx *player.Context, cmd cmd.Command, args []string) {
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
		*damage *= 1.5
	}
}

const scoreTagFormat = "\uE10C %d <red>|</red> %dms"

func (pr *Practice) HandleClientPacket(ctx *player.Context, pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.PlayerAuthInput:
		if p := instance.LookupPlayer(ctx.Val()); p != nil {
			if i := p.Instance(); i != instance.Nop {
				scoreTagTask(i, p)
				updateInputTask(pk, p)

			}
		}
	}
}

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

func updateInputTask(pk *packet.PlayerAuthInput, p *instance.Player) {
	if u, ok := p.User(); ok {
		currentMode := user.InputMode(pk.InputMode)

		if u.Dystopia().InputMode() != currentMode {
			u.Dystopia().SwitchInputMode(currentMode)
		}
	}
}

func (pr *Practice) HandleAttackEntity(ctx *player.Context, attacked world.Entity, force, height *float64, crit *bool) {
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

func welcomePlayer(p *player.Player) {
	_ = p.SetHeldSlot(4)

	welcomeTitle := title.New(text.Colourf("<red>Dystopia</red>"))
	welcomeTitle = welcomeTitle.WithFadeInDuration(time.Second * 2).WithDuration(time.Second).WithFadeOutDuration(time.Second)
	welcomeTitle = welcomeTitle.WithSubtitle(text.Colourf("<white>Welcome, <grey>%s</grey>!</white>", p.Name()))

	p.SendTitle(welcomeTitle)
}
