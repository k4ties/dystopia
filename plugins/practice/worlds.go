package practice

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/dystopia/embeddable"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/ffa/nodebuff"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/mw"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
)

type WorldConfig struct {
	Config struct {
		Name string

		Rotation [2]float64
		Spawn    [3]float64

		HeightThreshold instance.HeightThresholdConfig
	}
}

var (
	//go:embed lobby.json
	lobbyJson []byte
	//go:embed arena.json
	arenaJson []byte
)

var (
	lobbyConfig = embeddable.MustJSON[WorldConfig](lobbyJson)
	arenaConfig = embeddable.MustJSON[WorldConfig](arenaJson)
)

func setupWorldsManager(m *plugin.Manager, path string) {
	if err := mw.NewManager(m.Srv().World(), path, m.Logger()); err != nil {
		panic(err)
	}

	registerLobby(mw.M(), m)
	registerNodebuff(mw.M(), m)
}

func registerLobby(mn *mw.Manager, m *plugin.Manager) {
	name := lobbyConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, lobbyConfig.Config.Spawn); err != nil {
		panic(err)
	}

	hidden := []hud.Element{
		hud.Health, hud.Hunger,
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), lobbyConfig.Config.Rotation, lobbyConfig.Config.HeightThreshold, hidden...)
	instance.Register(name, in)
}

func registerNodebuff(mn *mw.Manager, m *plugin.Manager) {
	name := arenaConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("must create a world called arena")
	}

	if err := mn.SetSpawn(name, arenaConfig.Config.Spawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), arenaConfig.Config.Rotation, arenaConfig.Config.HeightThreshold)
	f := ffa.New(in.(*instance.Impl), nodebuff.Kit, ffa.Config{
		Name: "NoDebuff",
		Icon: "textures/items/potion_bottle_splash_heal.png",
	})

	instance.Register("nodebuff", f)
}
