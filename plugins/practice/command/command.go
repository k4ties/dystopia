package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

const owner = "i went tho"

type onlyPlayer struct{}

func (onlyPlayer) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}

type onlyOwner struct{}

func (onlyOwner) Allow(src cmd.Source) bool {
	p, ok := src.(*player.Player)
	return ok && p.Name() == owner
}

func systemMessage(o *cmd.Output, format string, args ...any) {
	o.Printf(text.Colourf("<red><b>>></b></red> %s", text.Colourf(format, args...)))
}

func p(s cmd.Source) *player.Player {
	return s.(*player.Player)
}

func inPl(s cmd.Source) *instance.Player {
	return instance.LookupPlayer(p(s))
}

func dead(s cmd.Source) bool {
	return p(s).GameMode() == world.GameModeSpectator
}
