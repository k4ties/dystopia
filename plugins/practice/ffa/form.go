package ffa

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Form struct{}

func (Form) Submit(s form.Submitter, pressed form.Button, tx *world.Tx) {
	if i, ok := instance.ByName(text.Clean(strings.Split(pressed.Text, "\n")[0])); ok {
		if f, ok := i.(*Instance); ok {
			if pl, in := instance.LookupPlayer(s.(*player.Player)); pl != nil && in != nil {
				f.Transfer(pl, tx)
			}
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

	sortButtons(buttons)
	return blank.WithButtons(buttons...).WithBody(text.Colourf("<white>Playing:</white> <grey>%d</grey>", len(players)))
}

func secondLine(ffa *Instance) string {
	if ffa.Closed() {
		return text.Colourf("<red>Closed</red>")
	}

	return text.Colourf("<grey>%d</grey> <dark-grey>Playing</dark-grey>", len(slices.Collect(ffa.Players())))
}

func isClosed(s string) bool {
	return text.Clean(getSecondLine(s)) == "Closed"
}

func getSecondLine(s string) string {
	return strings.Split(s, "\n")[1]
}

func mustPlayerCount(s string) int {
	c, err := strconv.Atoi(strings.Split(getSecondLine(text.Clean(s)), " ")[0])
	if err != nil {
		panic(err)
	}

	return c
}

func sortButtons(buttons []form.Button) {
	sort.Slice(buttons, func(i, j int) bool {
		if isClosed(buttons[i].Text) && isClosed(buttons[j].Text) {
			return len(buttons[i].Text) > len(buttons[j].Text)
		}
		if isClosed(buttons[i].Text) && !isClosed(buttons[j].Text) {
			return false
		}
		if !isClosed(buttons[i].Text) && isClosed(buttons[j].Text) {
			return true
		}
		if isClosed(buttons[i].Text) {
			return false
		}
		if isClosed(buttons[j].Text) {
			return true
		}

		if mustPlayerCount(buttons[i].Text) == mustPlayerCount(buttons[j].Text) {
			return len(buttons[i].Text) > len(buttons[j].Text)
		}

		return mustPlayerCount(buttons[i].Text) > mustPlayerCount(buttons[j].Text)
	})
}
