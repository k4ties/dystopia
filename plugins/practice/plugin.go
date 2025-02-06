package practice

import (
	_ "embed"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	_ "github.com/k4ties/dystopia/plugins/practice/command"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
)

//go:embed plugin.toml
var config []byte

func Plugin(l *handlers.Login, worldsPath string) plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		// ...
	})

	setupWorldsManager(m, worldsPath)
	return plugin.New(plugin.MustUnmarshalConfig(config), task, l, handlers.Practice{})
}

func LoginHandler(whitelisted bool, players ...string) plugin.PlayerHandler {
	return handlers.NewLoginHandler(whitelisted, players)
}
