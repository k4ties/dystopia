package ffa

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"slices"
	"strings"
)

type Form struct{}

func (Form) Submit(s form.Submitter, pressed form.Button, tx *world.Tx) {
	if i, ok := instance.ByName(strings.ToLower(pressed.Text)); ok {
		if f, ok := i.(*Instance); ok {
			f.Transfer(instance.LookupPlayer(s.(*player.Player)), tx)
		}
	}
}

func NewForm() form.Menu {
	blank := form.NewMenu(Form{}, "§g§r§fFFA")

	var buttons []form.Button
	var players []*instance.Player

	for _, inst := range instance.AllInstances() {
		if ffa, isFfa := inst.(*Instance); isFfa {
			nameFormat := "<white>%s</white>\n%s"

			for p := range ffa.Players() {
				players = append(players, p)
			}

			buttons = append(buttons, form.Button{
				Text:  text.Colourf(nameFormat, ffa.Name(), secondLine(ffa)),
				Image: ffa.Icon(),
			})
		}
	}

	return blank.WithButtons(buttons...).WithBody(text.Colourf("<white>Playing:</white> <grey>%d</grey>", len(players)))
}

func secondLine(ffa *Instance) string {
	if ffa.Closed() {
		return text.Colourf("<red>Closed</red>")
	}

	return text.Colourf("<grey>%d</grey> <red>Playing</red>", len(slices.Collect(ffa.Players())))
}
