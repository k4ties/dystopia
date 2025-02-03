package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type onlyPlayer struct{}

func (onlyPlayer) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}

func systemMessage(s cmd.Source, format string, args ...any) {
	if p, ok := s.(*player.Player); ok {
		p.Messagef(text.Colourf("<red><b>>></b></red> ")+format, args...)
	}
}
