package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/user"
	"strings"
)

type info struct {
	onlyOwner
	Player OnlinePlayer
}

func (c info) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	u, ok := user.Lookup(user.Name(c.Player))
	if !ok {
		o.Errorf("No player with name %s.", c.Player)
		return
	}

	format := []string{
		"InputMode: %s",
		"OS: %s",
		"Kills: %d",
		"Deaths: %d",
		"Rank: %s",
		"First join: %s",
		"IPS: %s",
		"DIDS: %s",
	}

	d := u.Dystopia()
	o.Printf(strings.Join(format, "\n"), d.InputMode().String(), d.OS().String(), d.Kills(), d.Deaths(), d.Rank().Name(), d.FirstJoin().String(), strings.Join(d.IPS(), ", "), strings.Join(d.DeviceIDs(), ", "))
}

func init() {
	cmd.Register(cmd.New("info", "shows info about player (session)", nil, info{}))
}

type OnlinePlayer string

func (OnlinePlayer) Type() string {
	return "player"
}

func (OnlinePlayer) Options(s cmd.Source) []string {
	var names []string

	for p := range plugin.M().Srv().Players(p(s).Tx()) {
		names = append(names, p.Name())
	}

	return names
}
