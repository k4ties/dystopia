package handlers

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Decliner struct{}

func (Decliner) HandleFoodLoss(ctx *player.Context, _ int, _ *int) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleFireExtinguish(ctx *player.Context, _ cube.Pos) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleStartBreak(ctx *player.Context, _ cube.Pos) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleBlockBreak(ctx *player.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleBlockPlace(ctx *player.Context, _ cube.Pos, _ world.Block) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleItemDamage(ctx *player.Context, _ item.Stack, _ int) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleItemPickup(ctx *player.Context, _ *item.Stack) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}

func (Decliner) HandleItemDrop(ctx *player.Context, _ item.Stack) {
	if ctx.Val().GameMode() != world.GameModeCreative {
		ctx.Cancel()
	}
}
