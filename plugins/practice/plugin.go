package practice

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/cmd"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/command"
	_ "github.com/k4ties/dystopia/plugins/practice/command"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
)

//go:embed plugin.toml
var config []byte

func Plugin(l *handlers.Login, worldsPath, databasePath string) plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		cmd.Register(cmd.New("ffa", "Requests FFA form", nil, command.FFA{}))
		cmd.Register(cmd.New("hub", "Teleports you to the lobby", nil, command.Hub{}))
		cmd.Register(cmd.New("rekit", "Resends kit if you're in ffa", nil, command.ReKit{}))
		cmd.Register(cmd.New("toggle", "Closes/opens specified FFA mode.", nil, command.Closer{}))
	})

	setupWorldsManager(m, worldsPath)
	return plugin.New(plugin.MustUnmarshalConfig(config), task, l, handlers.NewPractice(lobby.Instance()))
}
