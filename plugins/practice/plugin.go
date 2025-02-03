package practice

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/command"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/mw"
)

//go:embed plugin.toml
var config []byte

func Plugin(l *handlers.Login, worldsPath string) plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		cmd.Register(cmd.New("hub", "cmd.hub", []string{"lobby", "spawn"}, command.Hub{}))
	})

	setupWorldsManager(m, worldsPath)
	return plugin.New(plugin.MustUnmarshalConfig(config), task, l, handlers.NewInstance(lobby.Instance()), handlers.Practice{})
}

func LoginHandler(whitelisted bool, players ...string) plugin.PlayerHandler {
	return handlers.NewLoginHandler(whitelisted, players)
}

var lobbySpawn = mgl64.Vec3{0, -59, 0}

func setupWorldsManager(m *plugin.Manager, path string) {
	if err := mw.NewManager(m.Srv().World(), path, m.Logger()); err != nil {
		panic(err)
	}

	registerLobby(mw.M(), m)
}

func registerLobby(mn *mw.Manager, m *plugin.Manager) {
	w, ok := mn.World("lobby")
	if !ok {
		panic("must create a world called lobby")
	}

	if err := mn.SetSpawn("lobby", lobbySpawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger())
	instance.Register("lobby", in)
}
